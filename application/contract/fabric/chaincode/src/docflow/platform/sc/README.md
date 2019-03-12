# Smart Contracts



The smart contracts library is the main library. That library connects external 
calls with internal methods. It includes `SmartContract` object with 
`Init` and `Invoke` methods. `Invoke` method calls the method specified in 
the passed name with the arguments specified in passed arguments. It also checks 
if the user has the right to call this method.