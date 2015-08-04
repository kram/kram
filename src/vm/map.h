// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#ifndef VM_MAP_H
#define VM_MAP_H

class Value;

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

    typedef std::unordered_map<const char*, Value*, Hasher,  Equals<const char*> > map;
}

#endif