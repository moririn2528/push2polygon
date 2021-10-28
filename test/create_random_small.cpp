#include "testlib.h"
#include <bits/stdc++.h>

using namespace std;

const int N_MAX=1e5;
const int A_MAX=1e9,B_MAX=1e9;

void create_random(const int N, const int A, const int B, string name){
    int n=rnd.next(N)+1;
    int i;
    int a,b;
   cout<<n<<endl;
    for(i=0;i<n;i++){
        a=rnd.next(A)+1,b=rnd.next(B)+1;
       cout<<a<<" "<<b<<endl;
    }
}

void create_random_small(){
	int num;
    create_random(100,1000,40,"random_Small_"+to_string(num)+".in");
}

void create_random_large(int num){
    create_random(N_MAX,A_MAX,B_MAX,"random_large_"+to_string(num)+".in");
}

void create_random_testA(int num){
    create_random(N_MAX,A_MAX,1000,"testA_"+to_string(num)+".in");
}

void create_random_testB(int num){
    create_random(N_MAX,A_MAX,10,"testB_"+to_string(num)+".in");
}

int main(int argc, char* argv[]){
	registerGen(argc, argv, 1);
	create_random_small();
}
