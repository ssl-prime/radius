package service

import (
	"fmt"
	"net/http"
	"radius/api/util"

	"github.com/rightjoin/aqua"
)

//RepoInfo ...
type RepoInfo struct {
	aqua.RestService `prefix:"radius/" root:"/" version:"1"`
	getFrontend      aqua.GET  `url:"/index"`
	getOpenIssues    aqua.POST `url:"/getopenissues"`
	form             aqua.POST `url:"/form"`
	//getLE24hOpenIssues     aqua.GET  `url:"/onedayissues"`
	//getGT24hLE7dOpenIssues aqua.GET  `url:"/upto7dayissues"`
	//getGT7dOpenIssues      aqua.GET  `url:"olderthen7daysissues"`
	//openIssues             aqua.GET  `url:"/openissues"`
}

//GetOpenIssues ...
func (atdn *RepoInfo) GetOpenIssues(j aqua.Aide) (response interface{}, err error) {
	var (
		repoOwner string
		repoName  string
	)
	if repoOwner, repoName, err = util.ValidateGetOpenIssues(j); err == nil {
		response, err = util.GetOpenIssues(j, repoOwner, repoName)
	}
	return
}

//Form ...
func (atdn *RepoInfo) Form(j aqua.Aide) (response interface{}, err error) {
	j.LoadVars()
	var (
		repoOwner string
		repoName  string
	)
	if j.PostVar[`url`] != `` {
		if repoOwner, repoName, err = util.GetRepoOwnerAndName(j.PostVar[`url`]); err == nil {
			response, err = util.GetOpenIssues(j, repoOwner, repoName)
		} else {
			response = err.Error()
		}
	}
	return
}

//GetFrontend ...
func (atdn *RepoInfo) GetFrontend(j aqua.Aide) (response interface{}, err error) {
	fmt.Println(j.Request.Method)
	j.LoadVars()
	w := j.Response
	r := j.Request
	if j.Request.Method == "GET" {
		http.ServeFile(w, r, `/home/coinmark-003/go/src/radius/frontend/index.html`)
	} else {
		http.ServeFile(w, r, `/home/coinmark-003/go/src/radius/frontend/index.html`)
	}
	return "index.html", nil
}

// //GetLE24hOpenIssues ...
// func (atdn *RepoInfo) GetLE24hOpenIssues(j aqua.Aide) (response interface{}, err error) {
// 	response, err = util.Get24hIssues()
// 	return

// }

// //GetGT24hLE7dOpenIssues ...
// func (atdn *RepoInfo) GetGT24hLE7dOpenIssues(j aqua.Aide) (response interface{}, err error) {
// 	response, err = util.GetGT24hLTE7dIssues(j)
// 	return
// }

// //GetGT7dOpenIssues ...
// func (atdn *RepoInfo) GetGT7dOpenIssues(j aqua.Aide) (response interface{}, err error) {
// 	response, err = util.GetGT7dOpenIssues(j)
// 	return
// }

//OpenIssues
// func (atdn *RepoInfo) OpenIssues(j aqua.Aide) (response interface{}, err error) {
// 	response, err = util.OpenIssues(j)
// 	return
// }
