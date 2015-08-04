#include <unordered_map>
#include <iostream>

namespace Kram_Map {
    template <class _Tp>
    struct Equals : public std::binary_function<_Tp, _Tp, bool> {
        bool operator()(const _Tp& __x, const _Tp& __y) const  
        {
            return strcmp( __x, __y ) == 0;
        };
    };


    struct Hasher {
        // BKDR hash algorithm
        int operator()(const char * str)const
        {
            int seed = 131;
            int hash = 0;

            while(*str)
            {
                hash = (hash * seed) + (*str);
                str ++;
            };

            return hash & (0x7FFFFFFF);
        }
    };

    typedef std::unordered_map<const char*, unsigned int, Hasher,  Equals<const char*> > map;
}


int main(){
    Kram_Map::map location_map;

    char a[10] = "ab";
    location_map.insert(Kram_Map::map::value_type(a, 10));
    char b[10] = "abc";
    location_map.insert(Kram_Map::map::value_type(b, 20));

    char c[10] = "abc";
    location_map.insert(Kram_Map::map::value_type(c, 20));
    
    Kram_Map::map::iterator it;

    if ((it = location_map.find("abc")) != location_map.end())
    {
        std::cout << location_map["abc"];
    }

    return 0;
}