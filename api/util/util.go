package util

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"radius/api/model"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/rightjoin/aqua"
)

// ValidateGetOpenIssues ...
func ValidateGetOpenIssues(j aqua.Aide) (string, string, error) {
	j.LoadVars()
	var (
		urlBody   model.GetURL
		err       error
		repoOwner string
		repoName  string
	)
	if err = json.Unmarshal([]byte(j.Body), &urlBody); err == nil {
		if _, err = govalidator.ValidateStruct(urlBody); err == nil {
			repoOwner, repoName, err = GetRepoOwnerAndName(urlBody.URL)
		} else {
			err = errors.New("invalid required param : " + err.Error())
		}
	} else {
		err = errors.New("unmarshal error : " + err.Error())
	}

	return repoOwner, repoName, err
}

//GetRepoOwnerAndName ...
func GetRepoOwnerAndName(url string) (string, string, error) {
	var (
		repoOwner string
		repoName  string
		err       error
	)
	if spitURL := strings.Split(url, "https://github.com/"); len(spitURL) == 2 {
		if splitRepo := strings.Split(spitURL[1], "/"); len(splitRepo) == 2 {
			repoOwner = splitRepo[0]
			repoName = splitRepo[1]
		} else {
			err = errors.New("invalid repo details in url")
		}
	} else {
		err = errors.New("invalid url")
	}
	return repoOwner, repoName, err
}

//GetOpenIssues ...
func GetOpenIssues(j aqua.Aide, repoOwner, repoName string) (interface{}, error) {
	//github url for open issues
	//https://api.github.com/repos/originprotocol/origin/issues
	var (
		issues []model.IssuesResp
		err    error
		resp   interface{}
	)
	page := 1
	for {
		no := strconv.Itoa(page)
		issuesURL := `https://api.github.com/repos/` + repoOwner + `/` +
			repoName + `/issues?page=` + no +
			`&client_id=xxxxx&client_secret=yyyy`
		if issues, err = getIssues(issuesURL); err == nil {
			if len(issues) == 0 {
				break
			} else {
				err = processIssuesInfo(issues, repoOwner, repoName)
			}
		} else {
			break
		}
		page++
	}
	if err == nil {
		resp, err = OpenIssues(j, repoOwner, repoName)
	} else {
		resp = "free api limit exceeded"
	}
	return resp, err
}

func getIssues(issuesURL string) ([]model.IssuesResp, error) {
	var (
		err      error
		resp     *http.Response
		issues   []model.IssuesResp
		respBody []byte
	)
	if resp, err = http.Get(issuesURL); err == nil {
		if respBody, err = ioutil.ReadAll(resp.Body); err == nil {
			if err = json.Unmarshal(respBody, &issues); err != nil {
				err = errors.New("free allowed api call exceeded")
			}
		}
	} else {
		err = errors.New("request error: " + err.Error())
	}
	return issues, err
}

//processIssuesInfo ...
func processIssuesInfo(issues []model.IssuesResp, repoOwner, repoName string) error {
	var (
		issueSlice []model.TrackIssueInfo
		err        error
	)
	for i := 0; i < len(issues); i++ {
		var data model.TrackIssueInfo
		data.IssueURL = issues[i].URL
		data.OwnerName = repoOwner
		data.RepoName = repoName
		data.State = issues[i].State
		data.CreatedAt = (issues[i].CreatedAt).Format("2006-01-02 15:04:05")
		issueSlice = append(issueSlice, data)

	}
	err = UpdateOrInsertIssues(issueSlice)
	return err
}

//UpdateOrInsertIssues ...
func UpdateOrInsertIssues(issueSlice []model.TrackIssueInfo) error {
	//connect to db
	var (
		db   *sql.DB
		err  error
		stmt *sql.Stmt
	)
	if db, err = ConnectDB(); err == nil {
		sqlStr, vals := prepareIssuesData(issueSlice)
		upateOrInsert := `insert into repo_issues (issue_url, state, owner_name, 
			repo_name, created_at) values ` + sqlStr
		if stmt, err = db.Prepare(upateOrInsert); err == nil {
			stmt.Exec(vals...)
			stmt.Close()
			db.Close()
		} else {
			err = errors.New("update or insert error" + err.Error())
		}
	} else {
		err = errors.New("connection issue: " + err.Error())
	}
	return err
}

//prepareDataPhoto ...
func prepareIssuesData(issueSlice []model.TrackIssueInfo) (string, []interface{}) {
	vals := []interface{}{}
	sqlStr := ""
	for _, val := range issueSlice {
		sqlStr += `( ?, ?, ?, ?, ?),`
		vals = append(vals, val.IssueURL, val.State, val.OwnerName,
			val.RepoName, val.CreatedAt)
	}
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	return sqlStr, vals
}

//Get24hIssues ... will return issues older upto 24h
func Get24hIssues(j aqua.Aide, repoOwner, repoName string) (interface{}, error) {
	var (
		db    *sql.DB
		err   error
		count string
	)
	if db, err = ConnectDB(); err == nil {
		currentTimeStamp := getCurrentTime()
		openIssueCountQry := `SELECT count(*) from repo_issues
							WHERE created_at <= '` + currentTimeStamp + `' - INTERVAL 1 MINUTE
							AND created_at >'` + currentTimeStamp + `' - INTERVAL 1 Day + INTERVAL 1 MINUTE 
							AND owner_name = '` + repoOwner + `' AND repo_name = '` + repoName + `'`
		var results *sql.Rows
		if results, err = db.Query(openIssueCountQry); err == nil {
			for results.Next() {
				if err = results.Scan(&count); err != nil {
					log.Println("scan error :", err)
				}
			}
			db.Close()
		} else {
			err = errors.New("24h qry issue " + err.Error())
		}
	} else {
		err = errors.New("connection issue: " + err.Error())
	}
	return count, err
}

//GetGT24hLTE7dIssues  ... will return issues older then 24 h upto 7 days
func GetGT24hLTE7dIssues(j aqua.Aide, repoOwner, repoName string) (interface{}, error) {
	var (
		db    *sql.DB
		err   error
		count string
	)
	if db, err = ConnectDB(); err == nil {
		currentTimeStamp := getCurrentTime()
		openIssueCountQry := `SELECT count(*) from repo_issues
							WHERE created_at <= '` + currentTimeStamp + `' - INTERVAL 1 MINUTE -INTERVAL 1 DAY
							AND created_at >'` + currentTimeStamp + `' - INTERVAL 7 Day + INTERVAL 1 MINUTE 
							AND owner_name = '` + repoOwner + `' AND repo_name = '` + repoName + `'`
		var results *sql.Rows
		if results, err = db.Query(openIssueCountQry); err == nil {
			for results.Next() {
				if err = results.Scan(&count); err != nil {
					log.Println("scan error :", err)
				}
			}
			db.Close()
		} else {
			err = errors.New("24h qry issue " + err.Error())
		}
	} else {
		err = errors.New("connection issue: " + err.Error())
	}
	return count, nil
}

//GetGT7dOpenIssues ... will return issues older then 7 days
func GetGT7dOpenIssues(j aqua.Aide, repoOwner, repoName string) (interface{}, error) {
	var (
		db    *sql.DB
		err   error
		count string
	)
	if db, err = ConnectDB(); err == nil {
		currentTimeStamp := getCurrentTime()
		openIssueCountQry := `SELECT count(*) from repo_issues
							WHERE created_at <='` + currentTimeStamp + `' - INTERVAL 7 DAY  - INTERVAL 1 MINUTE
							AND owner_name = '` + repoOwner + `' AND repo_name = '` + repoName + `'`
		var results *sql.Rows
		if results, err = db.Query(openIssueCountQry); err == nil {
			for results.Next() {
				if err = results.Scan(&count); err != nil {
					log.Println("scan error :", err)
				}
			}
			db.Close()
		} else {
			err = errors.New("24h qry issue " + err.Error())
		}
	} else {
		err = errors.New("connection issue: " + err.Error())
	}
	return count, nil
}

func getCurrentTime() string {
	currentTime := time.Now()
	currentTimeStamp := (currentTime).Format("2006-01-02 15:04:05")
	return currentTimeStamp
}

//OpenIssues ...
func OpenIssues(j aqua.Aide, repoOwner, repoName string) (interface{}, error) {
	oneDayOpenIssue, _ := Get24hIssues(j, repoOwner, repoName)
	afterOneDayUpto7Days, _ := GetGT24hLTE7dIssues(j, repoOwner, repoName)
	olderThen7Days, _ := GetGT7dOpenIssues(j, repoOwner, repoName)
	oneDay, _ := strconv.Atoi(oneDayOpenIssue.(string))
	senvenDay, _ := strconv.Atoi(afterOneDayUpto7Days.(string))
	after7day, _ := strconv.Atoi(olderThen7Days.(string))
	totalOpenIssues := oneDay + senvenDay + after7day
	var issuesResp model.OpenIssues
	issuesResp.TotalOpenIssue = totalOpenIssues
	issuesResp.OneDay = oneDay
	issuesResp.Upto7Day = senvenDay
	issuesResp.OlderThen7Day = after7day
	return issuesResp, nil
}
