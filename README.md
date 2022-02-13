# Link Shortener for Ozon.FinTech
 
## Using docker:
Download or clone the project.

    docker build . --tag ozon_links
    docker run -p {your port}:15001 ozon_links:latest --[memory type]

Memory can be 2 types:

`--use-memory` - all data will be deleted after stopping the program.

`--use-db` - all data will be stored in postgreSql. 

If you use `--use-db`, then you need to edit Dockerfile's line 5 with your connection string.

In any case (except you want to test my application on port 15001), you should edit Dockerfile's line 3 and add your port. This base is used to return you correct values of short links.

## Testing API:
You can use Postman collection to check all endpoints:
1) POST `http://localhost:15001/` - to add new link for the server. 
2) GET `http://localhost:15001/` - to get long link to entered short link
3) GET `http://localhost:15001/{shortId}` - to be redirected into site, that is pointed by your short link.

I used REST API, so all inputs and outputs are done, using json. Input for all endpoints, except the third one are the same: 

    {
        "link": "your link"
    }
