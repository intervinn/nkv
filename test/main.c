#include <stdio.h>
#include <stdlib.h>
#include "kv.h"

#define TRY(expr) \
    if ((res = (expr)) != 0) goto fail; 

int main() {
    kv_client_t* c = kv_new("127.0.0.1", 8080);
    
    if (kv_connect(c) < 0) {
        perror(kv_err(c));
        return 1;
    }

    bool exists = false;
    char* val;
    int res;

    TRY(kv_exists(c, "foo", &exists));
    if (exists) {
        printf("alright, it exists\n");
    } else {
        printf("alright, it doesn't exist\n");
    }

    TRY(kv_put(c, "foo", "bar"));
    printf("okay we've put it\n");
    TRY(kv_get(c, "foo", &val));
    printf("alright we fetched it: %s\nnow time to delete it\n", val);
    TRY(kv_del(c, "foo"))

    free(val);
    kv_close(c);

    return 0;

    fail:
        printf("error: %s\n", kv_err(c));
        kv_close(c);
}