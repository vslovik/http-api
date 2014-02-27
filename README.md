## sCores Admin - Backend v1 for THE Football App

### Section management endpoints


**Create section**
<blockquote>POST /section/</blockquote>

Example:
```
{
    "name":"XXX",
    "published":0,
    "priority":0
}
```

**Read section**
<blockquote>GET /section/{sectionId}</blockquote>

**Update section**
<blockquote>PUT /section/{sectionId}</blockquote>

Example:
```
{
    "name":"xxx",
    "published":1,
    "priority":1
}
```

**Delete section**
<blockquote>DELETE /section/{sectionId}</blockquote>

**Add section translation**
<blockquote>POST /section/{sectionId}/translation/{language}</blockquote>

Example:
```
{
  "name": "xxx"
}
```

**Remove section translation**
<blockquote>DELETE /section/{sectionId}/translation/{language}</blockquote>

**List sections**
<blockquote>GET /section/</blockquote>

### Competition management endpoints

**Read competition**
<blockquote>GET /competition/{competitionId}</blockquote>

**Add competition translation**
<blockquote>POST /competition/{competitionId}/translation/{language}</blockquote>

Example:
```
{
  "name": "xxx"
}
```

**Remove competition translation**
<blockquote>DELETE /competition/{competitionId}/translation/{language}</blockquote>

**List competitions**
<blockquote>GET /competition/?page={page}</blockquote>
Number of competitions per page is set to 20.

**Add competition to section**
<blockquote>POST /section/{sectionId}/competition/{competitionId}</blockquote>

**Remove competition from section**
<blockquote>DELETE /section/{sectionId}/competition/{competitionId}</blockquote>

**Publish competition**
<blockquote>POST /competition/{competitionId}</blockquote>

**Hide competition**
<blockquote>DELETE /competition/{competitionId}</blockquote>



### Competition management endpoints

**Add competition to top list**
<blockquote>POST /top_competition/{countryId}/competition/<competition_id></blockquote>

**Remove competition from top list**
<blockquote>DELETE /top_competition/{countryId}/competition/<competition_id></blockquote>

**Read top competition list**
<blockquote>GET /top_competition/{countryId}</blockquote>

**Delete top competition list**
<blockquote>DELETE /top_competition/{countryId}</blockquote>


### Countries

**List countries**
<blockquote>GET /country/</blockquote>


### How to run

```
    go run main.go
```

 then send requests to port 8484
