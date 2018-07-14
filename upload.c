//
// Receive File Handler for multipart/form-data HTTP POST / HTML Form File Uploads
// Copyright (c) 1994-2017 Antoni Sawicki
//

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdarg.h>
#include <ctype.h>
#include <errno.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <unistd.h>

#include "cgic.h"

#define UPLOAD_DIR "/tmp"


void error(char *msg, ...) {
    va_list ap;
    char buff[1024]={0};

    va_start(ap, msg);
    vsnprintf(buff, sizeof(buff), msg, ap);
    va_end(ap);

    cgiHeaderContentType("text/plain");
    fprintf(cgiOut, "Error: %s\n%s [%d]\n", buff, strerror(errno), errno);

    exit(0);
}


int strip(char *str, int len, char *allow) {
    int n,a;
    int alpha=0, number=0, other=0;
    char *dst;

    if(!str || !strlen(str) || !allow || !strlen(allow))
        return -1;

    if(*allow == 'a') {
        alpha=1;
        allow++;
    }
    if(*allow == 'n') {
        number=1;
        allow++;
    }
    if(*allow == '!') {
        allow++;
    }

    if(strlen(allow))
        other=1;

    dst=str;

    for(n=0; n<len && *str!='\0'; n++, str++) {
        if(alpha && isalpha(*str))
            *(dst++)=*str;
        else if(number && isdigit(*str))
            *(dst++)=*str;
        else if(other)
            for(a=0; a<strlen(allow); a++) 
                if(*str==allow[a])
                    *(dst++)=*str;
    }

    *dst='\0';

    return 0;
}


int cgiMain(void) {
    cgiFilePtr input;
    FILE *output;
    char buff[8192]={0};
    char filename[128]={0};
    char fullfilename[1024]={0};
    int got=0;
    char *basename;
    char myurl[128]={0};
    struct stat mystat;

    if(cgiFormFileName("filename", filename, sizeof(filename)) == cgiFormSuccess) {

        if(cgiFormFileOpen("filename", &input) != cgiFormSuccess) 
                error("Unable to access uploaded file");

        basename=strrchr(filename, '/');
        if(!basename)
            basename=strrchr(filename, '\\');
                
        if(!basename)
            basename=filename;
        else
            basename++;

        strip(basename, sizeof(filename), "an_-.");
        snprintf(fullfilename, 1024, "%s/%s", UPLOAD_DIR, basename);

        output=fopen(fullfilename, "w");
        if(!output) 
                error("Unable to open file %s for writing", fullfilename);

        while(cgiFormFileRead(input, buff, sizeof(buff), &got) == cgiFormSuccess) 
            if(got)
                if(fwrite(buff, got, 1, output) != 1) 
                    error("While writing file %s [%s]", basename, strerror(errno));
        
        cgiFormFileClose(input);
        fclose(output);

        stat(fullfilename, &mystat);
        
        cgiHeaderContentType("text/plain");
        fprintf(cgiOut, "Received: %s, Size: %jd KB\n", basename, mystat.st_size/1024);
            
    }
    else {
        snprintf(myurl, sizeof(myurl), "http://%s%s", cgiServerName, cgiScriptName);
            
        if(strncmp(cgiUserAgent, "curl", 4)==0) {
            cgiHeaderContentType("text/plain");
            fprintf(cgiOut, "curl -F \"filename=@<yourfile.dat>\" %s\n", myurl);
        } else {
            cgiHeaderContentType("text/html");
            fprintf(cgiOut,
                "<HTML>\n"
                "  <BODY>\n"
                "Simple upload file via HTTP POST.<P>\n"
                "    <FORM ACTION=\"%s\" METHOD=\"POST\" ENCTYPE=\"multipart/form-data\">\n"
                "      <INPUT TYPE=\"file\" NAME=\"filename\">\n"
                "      <INPUT TYPE=\"submit\" NAME=\"Upload\" VALUE=\"Upload\">\n"
                "    </FORM><P>\n"
                "You can also use: <I>curl -F \"filename=@&lt;yourfile.dat&gt;\" %s</I><P>\n"
                "  </BODY>\n"
                "</HTML>\n", cgiScriptName, myurl);
        }

    }
    return 0;
}
