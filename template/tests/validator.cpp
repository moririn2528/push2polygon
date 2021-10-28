#include "testlib.h"
#include <bits/stdc++.h>

using namespace std;

constexpr int MIN_N = 1;
constexpr int MAX_N = 100000;
constexpr int MIN_AB = 1;
constexpr int MAX_AB = 1000000000;

int main() {
    registerValidation();

    int n = inf.readInt(MIN_N, MAX_N);
    inf.readEoln();

    for (int i = 0; i < n; i++) {
        int a = inf.readInt(MIN_AB, MAX_AB);
        inf.readSpace();
        int b = inf.readInt(MIN_AB, MAX_AB);
        inf.readEoln();
    }

    inf.readEof();

    return 0;
}
