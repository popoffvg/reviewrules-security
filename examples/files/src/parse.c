#include <stdio.h>
#include <string.h>

void set_label(const char *input) {
    char label[32];
    strcpy(label, input);
    puts(label);
}

void set_title(const char *input) {
    char title[32];
    snprintf(title, sizeof title, "%s", input);
    puts(title);
}
