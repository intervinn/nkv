#include <stdio.h>
#include "kv.h"

int main() {
    kv_client_t* c = kv_new("127.0.0.1", 8080);
    
    if (kv_connect(c) < 0) {
        perror(kv_err(c));
        return 1;
    }

    kv_exists(c, "foo", NULL);

    
    kv_close(c);
    return 0;
}