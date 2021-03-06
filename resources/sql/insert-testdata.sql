--  testdata
insert into monitor(monname,montype,monstatus,url,intrvl,exp_resp_code,timeout,retries) values("localhost","http","active","http://localhost:3090",12,200,1,4);
insert into monitor(monname,montype,monstatus,url,intrvl,exp_resp_code,timeout,retries) values("computerhok-http","http","active","http://www.computerhok.nl",10,404,7,1);
insert into monitor(monname,montype,monstatus,url,intrvl,exp_resp_code,timeout) values("computerhok-https","http","active","https://www.computerhok.nl",15,200,3);
insert into monitor(monname,montype,url,intrvl,retries) values("google","http","https://www.google.com/notthere",20,5);
insert into monitor(monname,montype,monstatus,url,intrvl) values("twitter","http","active","https://twitter.com/notthere",40);
insert into monitor(monname,montype,monstatus,url,intrvl) values("rabo-status","http","active","https://status.rabobank.nl/actuator/health",60);
insert into monitor(monname,montype,monstatus,url,intrvl) values("computerhok-wiki","http","active","http://www.computerhok.nl/wiki/Wiki.jsp?page=Main",70);

insert into chat(chatid, name) values(-235825137,"SchmutzigerHarry");
insert into chat(chatid, name) values(337345957,"AnnePup");

insert into username(name,role) values("SchmutzigerHarry","admin");
insert into username(name,role) values("AnnePup","reader");
