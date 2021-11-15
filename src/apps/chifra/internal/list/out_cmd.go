package listPkg

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

import (
	exportPkg "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/internal/export"
	"github.com/spf13/cobra"
)

var Options ListOptions

func RunList(cmd *cobra.Command, args []string) error {
	opts := Options

	exportPkg.Options.Appearances = true
	if opts.Count {
		exportPkg.Options.Count = true
	}
	if opts.FirstBlock > 0 {
		exportPkg.Options.FirstBlock = opts.FirstBlock
	}
	if opts.LastBlock > 0 {
		exportPkg.Options.LastBlock = opts.LastBlock
	}

	return exportPkg.RunExport(cmd, args)
}
