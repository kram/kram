// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include <fstream>
#include <iostream>
#include <string>

struct Output
{
	Output(std::ostream& os) : m_log(os.rdbuf()) { }

	std::streambuf* reset(std::ostream& os) 
	{
		return m_log.rdbuf(os.rdbuf());
	}

	template <typename T> friend Output& operator<<(Output& os, const T& t)
	{
		os.m_log << t;
		return os;
	}

	friend Output& operator<<(Output& os, std::ostream& ( *pf )(std::ostream&))
	{
		os.m_log << pf;
		return os;
	}

	private:
		std::ostream m_log;
};