#pragma once
/*-------------------------------------------------------------------------------------------
 * qblocks - fast, easily-accessible, fully-decentralized data from blockchains
 * copyright (c) 2016, 2021 TrueBlocks, LLC (http://trueblocks.io)
 *
 * This program is free software: you may redistribute it and/or modify it under the terms
 * of the GNU General Public License as published by the Free Software Foundation, either
 * version 3 of the License, or (at your option) any later version. This program is
 * distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even
 * the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details. You should have received a copy of the GNU General
 * Public License along with this program. If not, see http://www.gnu.org/licenses/.
 *-------------------------------------------------------------------------------------------*/
#include "acctlib.h"

#include "classes/cachebase.h"
#include "classes/cache.h"

#include "classes/chaincache.h"
#include "classes/abicacheitem.h"
#include "classes/abicache.h"
#include "classes/collectioncacheitem.h"
#include "classes/collectioncache.h"
#include "classes/monitorcacheitem.h"
#include "classes/monitorcache.h"
#include "classes/indexcacheitem.h"
#include "classes/indexcache.h"
#include "classes/pricecacheitem.h"
#include "classes/pricecache.h"
#include "classes/namecache.h"
#include "classes/slurpcache.h"

#include "classes/status.h"
#include "classes/statusterse.h"

namespace qblocks {
extern bool countFiles(const string_q& path, void* data);
extern bool countFilesInCache(const string_q& path, void* data);
}  // namespace qblocks

using namespace qblocks;  // NOLINT
