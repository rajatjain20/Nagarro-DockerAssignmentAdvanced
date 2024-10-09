Project Overview:

	I have built a multiservice application, where there will be three services - 

	1. Data input: 
		This is a Golang based web application, through which, user can add text data (.text file) in a shared directory (using webapi) from where 	the webapp will pick that data

	2. Webapp:
		This is a Golang based web application, which will access the data stored in the shared directory from "DataInputApp" and store that data 	into Database (using webapi). User will be able to view that data stored in DB through a webapi.

	3. Database:
		This is a MSSQL database. In this, we have 'FILE2DB' database and a table named "FILE2DBDATA". Data will be stored in this table only by 	webapp.


Setting up volumes, networks, images and containers -

 
	1. Create Volume for database persistence :
		> docker volume create mssqldb_volume

	2. Create Volume for sharing data between dataInputapp and webapp :
		> docker volume create data_volume

	3. Create Network on which ‘webapp’ container can communicate with ‘mssqldb’ container:
		> docker network create -d bridge webapp_network

	4. Create images and run containers – 

		i. db image and container: (go to db directory in project on terminal)

			> docker build -t mssqldb:v0 .
			>  docker run -e MSSQL_SA_PASSWORD=admin@123 -p 1433:1433 -v mssqldb_volume:/var/opt/mssql --name mssqldb -d mssqldb:v0

	   	   Connect DB container to network for communication with webapp:

			> docker network connect --alias mssqldb webapp_network mssqldb
			> docker inspect mssqldb

		ii. datainputapp image and container: (go to datainputapp directory in project on terminal)
		
			> docker build -t datainputapp:v0 .
			> docker run --env-file ./config/.env-dev -v data_volume:/go/inputdata/ -p 3400:3400 --name datainputapp datainputapp:v0

		iii. webapp image and container: (go to webapp directory in project on terminal)

			> docker build -t webapp:v0 .
			> docker run --env-file ./config/.env-dev --volumes-from datainputapp --network webapp_network -p 4000:4000 --name webapp webapp:v0


So now all three containers are running and they are connected through network and required volumes are mounted.
Let's execute and test the application.

Below are the webapi:

	1. datainputapp: (running on port-3400)

		i. http://localhost:3400/
		ii. http://localhost:3400/writedata?filename=file1.txt&data=Hello!! This is file1.txt

	2. webapp: (running on port-4000)

		i.  http://localhost:4000/
		ii. http://localhost:4000/processdata
		iii.http://localhost:4000/showdbdata

To Check data persistence:

	We will stop ‘mssqldb’ container and remove it and then create and run it again.
		> docker stop mssqldb
		> docker container rm mssqldb
		> docker ps -a

	Now run the container again and also remember that we have to add network to it as we did before:
		> docker run -e MSSQL_SA_PASSWORD=admin@123 -p 1433:1433 -v mssqldb_volume:/var/opt/mssql --name mssqldb -d mssqldb:v0
		> docker network connect --alias mssqldb webapp_network mssqldb

	Now if we run below url, the data should persist:
		http://localhost:4000/showdbdata



Docker compose:

	There is one ‘docker-compose.yml’ file has been created in the project.
	
	On the terminal, go to File2DB directory and run below command:
		> docker compose up -d
	
	Lets look at the images:
		> docker images
	
	Container for this compose file:
		> docker ps

	> docker compose down


Push images to docker registry (Docker Hub):

Commands:
> docker login
> docker tag webapp:v0 rajatjain20/webapp:v0
> docker push rajatjain20/webapp:v0

> docker tag datainputapp:v0 rajatjain20/datainputapp:v0
> docker push rajatjain20/datainputapp:v0

> docker tag mssqldb:v0 rajatjain20/mssqldb:v0
> docker push rajatjain20/mssqldb:v0


CI/CD pipeline:

I have used GitLab for CI/CD of this application. I have created a file ".gitlab-ci.yml", which builds the images of the applications and saves them into GitLab registry.

Note: Please refer "Project Overview and Screenshots.docx" for more details and screenshots.