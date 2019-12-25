package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	mapset "github.com/deckarep/golang-set"
	_ "github.com/go-sql-driver/mysql"
)

//数据库连接信息
const (
	USERNAME = "root"
	PASSWORD = "741235896"
	NETWORK  = "tcp"
	SERVER   = "localhost"
	PORT     = 13306
	DATABASE = "test"
)

func initDB() *sql.DB {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, err := sql.Open("mysql", conn)
	checkErr(err)
	if db != nil {
		fmt.Println("DB Connected!")
	}
	db.SetConnMaxLifetime(100 * time.Second) //最大连接周期，超时的连接就close
	db.SetMaxOpenConns(100)
	return db
}

type Entity struct {
	id         int
	name       string
	lat        float64
	lng        float64
	sn         string
	ip         string
	cameraIP   string
	cameraName string
	areaID     int64
}

// predeclare err
func parseCsv(file string, areaID int64) (map[string]*Entity, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close() // this needs to be after the err check

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	devices := make(map[string]*Entity, len(lines))
	var lat, lng float64
	var name, sn, ip, cameraIP, cameraName string
	for idx, line := range lines[1:] {
		name = line[0]
		ip = line[3]
		sn = line[4]
		cameraIP = line[5]
		cameraName = line[6]
		if lat, err = strconv.ParseFloat(line[2], 64); err != nil {
			return nil, err
		}
		if lng, err = strconv.ParseFloat(line[1], 64); err != nil {
			return nil, err
		}
		devices[sn] = &Entity{id: idx, lat: lat, lng: lng, name: name, sn: sn, ip: ip, cameraIP: cameraIP, cameraName: cameraName, areaID: areaID}
	}
	return devices, nil
}

func readAreas(db *sql.DB) {
	//查询数据
	rows, err := db.Query("SELECT id,name FROM dvc_area where delflag=0")
	checkErr(err)

	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		checkErr(err)
		fmt.Println(id, name)
	}
}

func readSn(db *sql.DB) mapset.Set {
	//查询数据
	rows, err := db.Query("SELECT sn FROM dvc_device where delflag=0")
	checkErr(err)

	snList := mapset.NewSet()
	for rows.Next() {
		var sn string
		err = rows.Scan(&sn)
		checkErr(err)
		snList.Add(sn)
	}

	return snList
}

func outputCsv(data map[string]*Entity) error {

}
func main() {
	var csvFile string
	flag.StringVar(&csvFile, "path", "data.csv", "data.csv")

	flag.Parse()
	fmt.Println("========== 选择要导入的区域编号 ==========")
	db := initDB()
	readAreas(db)
	existsSnList := readSn(db)
	fmt.Print(">>>>>>  输入要导入的区域编号，回车确认：")
	var areaID string
	fmt.Scanln(&areaID)
	intAreaID, err := strconv.ParseInt(areaID, 10, 0)
	devices, err := parseCsv(csvFile, intAreaID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("待导入数据量:", len(devices))
	tx, err := db.Begin() // 创建tx对象
	errFlag := false
	ignoreDevices := make([]*Entity, 10)

	total := len(devices)
	succ := 0
	ignored := 0
	for _, device := range devices {
		if existsSnList.Contains(device.sn) {
			// exists
			ignoreDevices = append(ignoreDevices, device)
			ignored++
			continue
		}
		// insert rows with transaction
		if _, err = tx.Exec("INSERT INTO dvc_device(name, sn, ip, gps_lat, gps_lng, area_id ) VALUES(?, ?, ?, ?, ?, ?);",
			device.name, device.sn, device.ip, device.lat, device.lng, device.areaID); err != nil {
			fmt.Println(err)
			errFlag = true
			break
		}
		if _, err = tx.Exec("INSERT INTO dvc_camera(name, device_sn, ip ) VALUES(?, ?, ?);",
			device.cameraName, device.sn, device.cameraIP); err != nil {
			fmt.Println(err)
			errFlag = true
			break
		}
		succ++
	}
	if errFlag {
		tx.Rollback()
		fmt.Println("导入失败!")
	} else {
		tx.Commit()
		fmt.Println("导入成功!")
		fmt.Println("共计:", total, ", 成功导入: ", succ, ", 忽略的重复项目: ", ignored)

	}

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
