---
title: Language Codeblock Test
author: John Doe
date: 2024-02-23
scripts:
collections: ["test", "posts"]
tags:
  - rust
  - C
---

## struct in rust

```rust
struct Vector2D {
    x: f64,
    y: f64,
}
```

## Brainfuck interpreter yanked from reddit

```C
s[999],*r=s,*d,c;main(a,b){char*v=1[d=b];for(;c=*v++%93;)for(b=c%7?a&&(c&17?c&1?
(*r-=c-44):(r+=c-61):c&2?putchar(*r):(*r=getchar()),0):v;b&&c|a**r;v=d)main(!c,&
b-1);d=v;}
```

## Entire `yake` sourcecode!

```C
#include <stdio.h>

int opr(char *v);
void ex(int *p, char *syn);

int main(int ac, char *av[]) {
    int p = 0;
    return (ac < 2) ? (printf("synerr\n"),1) : (ex(&p, av[1]), 0);
}

void ex(int *v, char *syn){
    int *p;
    int st = 0; // 0 for value control, 1 for pointer control
    while ((int)*syn!='\0'){
        switch ((int)*syn) {
            case 46: printf("&v : %d | %p -> %d | m = %d\n", *v, p, *p, st); break;
            case 43: if(st!=1){++*v;}else{++p;} break;
            case 45: if(st!=1){--*v;}else{--p;} break;
            case 60: if(st!=1){*v = *v<<1;}else{printf("synerr");return;} break;
            case 62: if(st!=1){*v = *v>>1;}else{printf("synerr");return;} break;
            case 64: if(st == 1){st = 0;*p=*v;}else{st=1;} break;
            case 38: if(st == 1){st = 0;*v=*p;}else{st=1;}break;// & is used to get value

            // operations, st == 1, for subsequent value for operation type
            case 37: if(st==1){
                syn++;
                switch((int)*syn) {
                    case 45: *v-=*p; break;
                    case 43: *v+=*p; break;
                }
                st = 0; // auto switch back to st=0
            } break;

            case 58: // ascii to string
                     // length of string to be printed is
                     // to be stored in v
                     for (int i=0;i<*v;++i){printf("%c",*p);++p;}

                     printf("\n");
                     break;

            default: printf("synerr\n"); return;
        }
        syn++;
    }
    return;
}
```
