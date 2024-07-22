# A snippet saving site under the name 'snippet box'
- coded using go
    - used the standard http lib, with the help of `julienschmidt/httprouter`
    - session handling done using `alexedwards/scs`
- Database: mysql

## Models and attributes
### Snippet
```go
ID // id of every snippet, which is unique number to diffrentiate between snippets even if they have the same title and content
Title // title of every snippet, is used to show the snippet in the home page 
Content // content of every snippet
Created // time which the snippet was created and is shown in the snippet view page
Expires // time which the snippet will expire at and is also shown in the snippet view page, users can set the expiration time while creating a snippet
```

### User
```go
ID // unique id for every user
Name // name of the user
Email // unique email of user, used in login
HashedPassword // hashed password to achieve security principles if the database was comprimised
Created // this date which every user account was created
```


