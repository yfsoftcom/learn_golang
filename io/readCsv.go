package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"text/template"
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
	TEMPLATE = `<head><style>td{padding: 3px;}</style></head><h2>导入结果</h2><p>导入时间:{{.now }}</p><p>共计: {{ .total }} 条，成功： {{ .success }} 条</p>
	<h3>可能重复的SN数据</h3>
<table border="1" cellspacing="0" cellpadding="0">
<thead><tr><th>名称</th><th>经度</th><th>纬度</th><th>ip</th><th>sn</th><th>监控ip</th><th>监控名称</th></tr></thead>
<tbody>
{{ range .devices }}
<tr ><td>{{ .Name}}</td><td>{{ .Lng}}</td><td>{{ .Lat}}</td><td>{{ .IP}}</td><td>{{ .Sn}}</td><td>{{ .CameraIP}}</td><td>{{ .CameraName}}</td></tr>
{{ end }}
</tbody></table>`
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
	ID         int
	Name       string
	Lat        float64
	Lng        float64
	Sn         string
	IP         string
	CameraIP   string
	CameraName string
	AreaID     int64
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
		devices[sn] = &Entity{ID: idx, Lat: lat, Lng: lng, Name: name, Sn: sn, IP: ip, CameraIP: cameraIP, CameraName: cameraName, AreaID: areaID}
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

func output(total, succ int, data []Entity) error {
	t, err := template.New("output").Parse(TEMPLATE)
	checkErr(err)

	now := time.Now().Format("2006-01-02 15:04:05")
	f, err := os.Create("import_" + now + ".html")
	defer f.Close()
	checkErr(err)
	return t.Execute(f, map[string]interface{}{
		"now":     now,
		"total":   total,
		"success": succ,
		"devices": data,
	})
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
	ignoreDevices := make([]Entity, 0)

	total := len(devices)
	succ := 0
	ignored := 0

	timeUnix := time.Now().Unix()
	for _, device := range devices {
		if existsSnList.Contains(device.Sn) {
			// exists
			ignoreDevices = append(ignoreDevices, *device)
			ignored++
			continue
		}
		// insert rows with transaction
		if _, err = tx.Exec("INSERT INTO dvc_device(name, sn, ip, gps_lat, gps_lng, area_id, createAt, updateAt ) VALUES(?, ?, ?, ?, ?, ?, ?, ?);",
			device.Name, device.Sn, device.IP, device.Lat, device.Lng, device.AreaID, timeUnix, timeUnix); err != nil {
			fmt.Println(err)
			errFlag = true
			break
		}
		if _, err = tx.Exec("INSERT INTO dvc_camera(name, device_sn, ip, qa_length, model, setup_date, createAt, updateAt ) VALUES(?, ?, ?, ?, ?, ?, ?, ?);",
			device.CameraName, device.Sn, device.CameraIP, 24, "kakou", timeUnix, timeUnix, timeUnix); err != nil {
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
		fmt.Println("共计:", total, "\n成功导入: ", succ, "\n忽略的重复项目: ", ignored)
		err := output(total, succ, ignoreDevices)
		checkErr(err)
	}

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
