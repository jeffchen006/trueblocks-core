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
 * Parts of this file were generated with makeClass --options. Edit only those parts of
 * the code outside of the BEG_CODE/END_CODE sections
 */
#include "options.h"

//---------------------------------------------------------------------------------------------------
bool COptions::writeOpenApiFile(void) {
    if (!openapi)
        return true;

    LOG_INFO(cYellow, "handling openapi file...", cOff);
    counter = CCounter();  // reset

    map<string_q, string_q> converts;
    map<string_q, string_q> pkgs;
    for (auto ep : endpointArray) {
        CCommandOptionArray params;
        for (auto option : routeOptionArray)
            if (option.api_route == ep.api_route && option.isChifraRoute(false))
                params.push_back(option);
        ep.params = &params;

        string_q productions = getProductions(ep);
        string_q exampleFn = getDocsPathTemplates("api/examples/" + ep.api_route + ".txt");

        chifraCmdStream << ep.toChifraCmd() << endl;
        pairMapStream << ep.toPairMap() << endl;
        apiTagStream << ep.toApiTag();
        goCallStream << ep.toGoCall();
        goRouteStream << ep.toGoRoute();
        apiPathStream << ep.toApiPath(productions, exampleFn);

        string_q pkg = ep.toGoPackage();
        string_q nick = nextTokenClear(pkg, ' ');
        pkgs[nick] = substitute(pkg, "\n", "");

        if (isApiRoute(ep.api_route) && ep.params) {
            for (auto p : *((CCommandOptionArray*)ep.params))
                if (contains(p.longName, "_"))
                    converts[toCamelCase(p.longName)] = p.longName;
        }

        counter.cmdCount += params.size();
        counter.routeCount++;
    }
    for (auto pkg : pkgs)
        goPkgStream << pkg.first << " " << pkg.second << endl;
    goPkgStream << "\tconfig \"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config\"" << endl;

    converts["logLevel"] = "log_level";
    converts["noHeader"] = "no_header";
    converts["toFile"] = "to_file";
    goConvertStream << "\tswitch in {" << endl;
    for (auto item : converts) {
        goConvertStream << "\tcase \"" << item.first << "\":\n\t\treturn \"" << item.second << "\"" << endl;
    }
    goConvertStream << "\t}" << endl;

    writeCodeOut(this, getPathToSource("apps/chifra/server/routes.go"));
    writeCodeOut(this, getPathToSource("apps/chifra/server/convert_params.go"));

    if (getEnvStr("GENERATE_YAML") != "false") {
        writeCodeOut(this, getDocsPathContent("api/openapi.yaml"));
    }
    if (getEnvStr("GENERATE_SDK") == "true") {
        ostringstream tsos;
#define explorerPath string_q("/Users/jrush/Development/trueblocks-explorer/")
#define coreDocsPath string_q("/Users/jrush/Development/trueblocks-core/docs/content/api/openapi.yaml")
        tsos << "cd " << explorerPath << " ; ";
        tsos << "URL=" << coreDocsPath << " yarn generate:sdk";
        if (system(tsos.str().c_str()) != 0) {
        }
    }

    LOG_INFO(cYellow, "makeClass --openapi", cOff, " processed ", counter.routeCount, "/", counter.cmdCount,
             " routes/cmds ", " (changed ", counter.nProcessed, ").", string_q(40, ' '));

    return true;
}

namespace qblocks {
extern bool isApiRoute(const string_q& route);
}

//---------------------------------------------------------------------------------------------------
string_q COptions::getProductions(const CCommandOption& ep) {
    if (!isApiRoute(ep.api_route) || contains(ep.api_route, "explore"))
        return "";
    CStringArray productions;
    string_q descr;
    for (auto model : dataModels) {
        if (contains(model.doc_producer, ep.api_route)) {
            productions.push_back(model.doc_api);
            if (descr.empty())
                descr = model.doc_descr;
        }
    }

    string_q prods;
    for (auto p : productions) {
        replace(p, "cachePtr", "cache");
        prods += "$ref: \"#/components/schemas/" + p + "\"\n";
    }
    if (productions.size() > 1) {
        prods = "oneOf:\n" + prods;
        replaceAll(prods, "$ref:", "      - $ref:");
    }

    return trim(prods, '\n') + "|" + descr;
}
