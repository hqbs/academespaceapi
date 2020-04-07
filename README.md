# dtaback

Version 1.0!

# Basic How To Use 
## Get request the API Example Below

### Sample Query
http://APIENDPOINTIP/graphql?query={validateUserToken(email:"STRING",currenttoken:"STRING")}

### Sample Mutation 
http://APIENDPOINTIP/graphql?query=mutation+_{createUser(fname:"STRING",lname:"STRING",email:"STRING",phonenumber:"STRING",type:"STRING",password:"STRING",discordid:"STRING"){Success Errors Token}}

## Queries and Return Types

userExists 
 - Arguments 
   - email
     - Type: String
 - Return
   - {} (empty return)
     - Type: Boolean

validateUserToken
 - Arguments
   - email
     - Type: String
   - currenttoken
     - Type: String
 - Return
   - {} (empty return)
     - Type: Boolean
  
login
 - Arguments
   - email
     - Type: String
   - password
     - Type: String
   - token
     - Type: String
 - Return
   - success
     - Type: Boolean
     - *Denotes any errors - false == errors*
   - errors
     - Type: []String
     - *String array of all errors*
   - token
     - Type: String
     - *Returns a valid token only when success is true*

## Mutations and Return Types

createUser
 - Arguments
   - fname
     - Type: String
   - lname
     - Type: String
   - email
     - Type: String
   - phonenumber
     - Type: String
   - type
     - Type: String 
   - password
     - Type: String
   - discordid
     - Type: String
 - Return
   - success
     - Type: Boolean
     - *Denotes any errors - false == errors*
   - errors
     - Type: []String
     - *String array of all errors*
   - token
     - Type: String
     - *Returns a valid token only when success is true*



