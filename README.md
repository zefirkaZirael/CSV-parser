ReadLine: 
This function reads a new line from a CSV file, returning the pointer to the line with the terminator removed. 
If the line has a missing or extra quote, it returns an empty string and an ErrQuote error. The function handles various 
line terminations (\r, \n, \r\n, or EOF) and can be called in a loop to sequentially read each line from the file. The
returned line does not include the \n at the end.

GetField:
This function retrieves the nth field from the last line read by ReadLine. If the field is out of range or if n is negative,
it returns an ErrFieldCount error. Fields are separated by commas, and any quotes around the fields are removed. This 
function handles an arbitrary number of fields with any length.

GetNumberOfFields: 
This function returns the number of fields in the last line read by ReadLine. This should only be called after ReadLine has
been invoked.
