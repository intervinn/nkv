#ifndef KV_H
#define KV_H

#include <stdbool.h>

typedef struct {
    int len;
    char** fragments;
} kv_message_t;

typedef struct _kv_client_t kv_client_t;

kv_client_t *kv_new(const char *addr, int port);
int kv_connect(kv_client_t *c);
int kv_close(kv_client_t *c);
const char *kv_err(kv_client_t *c);

int kv_exists(kv_client_t *c, char *key, bool *value);
int kv_put(kv_client_t *c, char *key, char *value);
int kv_get(kv_client_t *c, char *key, char **value);
int kv_del(kv_client_t *c, char *key);

#endif