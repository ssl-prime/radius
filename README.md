# radius
tracking repository open issues

to the run this repo need to install golang 
set go path 
clone this repo in src folder
basic changes file location , database connection
 create database radius :
create table `repo_issues` (
  `id` int(11) auto_increment primary key not null,
  `issue_url` varchar(256) not null unique key,
  `state` enum('open','closed') default 'open',
  `owner_name` varchar(256) not null,
  `repo_name` varchar(256) not null,
  `created_at`  datetime ,
  `updated_at` datetime,
  `closed_at` datetime
  );
move to launch 
then use this `go get` 
this will install all required package 
now `go run main.go`


basic url(localhost) for operation = `http://localhost:8090/radius/v1/index` 
basic url aws = `http://52.91.74.120:8090/radius/v1/index`

