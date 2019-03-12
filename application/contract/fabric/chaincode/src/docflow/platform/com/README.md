# Communication

The communication library enables you to send response to the server 
with Error or Success responses and write this response to logs. You also 
can use this methods for logs only.
 
There are three fields in response:

* `status` of response
* error or success `message`
* data of `payload`

Error response have only code and message. Success response can have 
code with message or code with data.



## Examples

Example of success response with data:

```
return com.SuccessPayloadResponse([]byte(result))
```

Example of success message for logs:

```
// Retrieve the requested Smart Contract function and arguments
function, args := APIstub.GetFunctionAndParameters()

com.InfoLogMsg("Invoke of " + function)
```


Example of error message when unmarshal:


```		
statement := statement.Statement{}

err := json.Unmarshal([]byte(value), &statement)
if err != nil {
return com.UnmarshalError(string(value))
}
```
