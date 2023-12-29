# Naming Convention
* Prefix keys- is put when the type is a slice with the length of keyCount.
* Prefix k- for local variables of keys-.
    * Put suffix 'List' when suffix 's' is not available.

* Index variable ends with -i suffix.
* It is fine to put plural suffix (-s) in the middle if putting the suffix at the last is not proper (e.g., adjective)
* Suffix -t for local variable name of -Type

* Avoid using abbreviation in struct name and field name 
unless the name is explicitly supposed to be expressed in abbreviated form.  

* Local variables are encouraged to be written in abbreviated form, which is up to 3 letters. 
    * field name: sprites, anims
    * local name: s, a

# Keyword
ratio: two quantities with the same units
rate: two quantities with the different units

# Code
* When there are local variable and receiver field which points identical value, use local variable: prefer to use shorter one. 

* The convention to organize types and structs in a file by defining the dependencies first and the types that utilize those dependencies later.
    * However, if a small function is used only in a single large function, then put the small at the below of the large.

* If you find any helpful information, commenting is preferred. Grammar check is not necessary at this stage. 

## Field order
### V1
// Order of fields: logic -> drawing.
// Try to keep the order of initializing consistent with the order of fields.
