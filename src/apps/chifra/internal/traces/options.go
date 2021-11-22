package tracesPkg

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
/*
 * This file was auto generated with makeClass --gocmds. DO NOT EDIT.
 */

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/internal/globals"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/validate"
)

type TracesOptions struct {
	Transactions []string
	Articulate   bool
	Filter       string
	Statediff    bool
	Count        bool
	SkipDdos     bool
	Max          uint64
	Globals      globals.GlobalOptions
	BadFlag      error
}

func (opts *TracesOptions) TestLog() {
	logger.TestLog(len(opts.Transactions) > 0, "Transactions: ", opts.Transactions)
	logger.TestLog(opts.Articulate, "Articulate: ", opts.Articulate)
	logger.TestLog(len(opts.Filter) > 0, "Filter: ", opts.Filter)
	logger.TestLog(opts.Statediff, "Statediff: ", opts.Statediff)
	logger.TestLog(opts.Count, "Count: ", opts.Count)
	logger.TestLog(opts.SkipDdos, "SkipDdos: ", opts.SkipDdos)
	logger.TestLog(opts.Max != 250, "Max: ", opts.Max)
	opts.Globals.TestLog()
}

func (opts *TracesOptions) ToCmdLine() string {
	options := ""
	if opts.Articulate {
		options += " --articulate"
	}
	if len(opts.Filter) > 0 {
		options += " --filter " + opts.Filter
	}
	if opts.Statediff {
		options += " --statediff"
	}
	if opts.Count {
		options += " --count"
	}
	if opts.SkipDdos {
		options += " --skip_ddos"
	}
	if opts.Max != 250 {
		options += (" --max " + fmt.Sprintf("%d", opts.Max))
	}
	options += " " + strings.Join(opts.Transactions, " ")
	options += fmt.Sprintf("%s", "") // silence go compiler for auto gen
	return options
}

func FromRequest(w http.ResponseWriter, r *http.Request) *TracesOptions {
	opts := &TracesOptions{}
	opts.Max = 250
	for key, value := range r.URL.Query() {
		switch key {
		case "transactions":
			for _, val := range value {
				s := strings.Split(val, " ") // may contain space separated items
				opts.Transactions = append(opts.Transactions, s...)
			}
		case "articulate":
			opts.Articulate = true
		case "filter":
			opts.Filter = value[0]
		case "statediff":
			opts.Statediff = true
		case "count":
			opts.Count = true
		case "skip_ddos":
			opts.SkipDdos = true
		case "max":
			opts.Max = globals.ToUint64(value[0])
		default:
			if !globals.IsGlobalOption(key) {
				opts.BadFlag = validate.Usage("Invalid key ({0}) in {1} route.", key, "traces")
				return opts
			}
		}
	}
	opts.Globals = *globals.FromRequest(w, r)

	return opts
}

var Options TracesOptions

func TracesFinishParse(args []string) *TracesOptions {
	// EXISTING_CODE
	Options.Transactions = args
	// EXISTING_CODE
	return &Options
}
