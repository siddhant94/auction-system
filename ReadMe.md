### Auction System

##### Steps to run the app
- Clone / Download  the Repo.
- Install Docker
- On Linux(Debian), run  `docker build -t auction-system .`  
 from the root directory of project.
 - On successful build, run `docker images` to get a list of available images
- Finally run `docker run -p 5000:8081 -it auction-system` to run the app on localhost, port 5000.