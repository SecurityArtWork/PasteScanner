# PasteScanner
Pastescanner is a (Golang)tool to monitor websites content for relevant information

The programme shows nothing on terminal, it generates a file in ./pastes for each paste without expiration time, and in ./pastes/temp with certain time to live pastes.

In order to configure it, edit paste.conf. It has two highly descriptive tags, you can add as much keywords as you want so it can filter pastes by them. It is also possible to erase/comment lines about different pastesites so they are not included in the search engine. It is not possible to add new pastesites without editing the code further, however it is quite easy as it's just writing a few lines of code in main() and find() and adding it in .conf. You're welcome to contribute!