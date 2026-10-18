import {Box, Link} from "@mui/material"
import {Fragment} from "react"

import {SxPropsMap} from "../../helper/HelperType"
import {IvoryLinks} from "../../helper/HelperUtils"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", fontFamily: "monospace", margin: "5px 20px 5px", gap: "5px"},
    links: {display: "flex", flexWrap: "wrap", justifyContent: "center", rowGap: 1, columnGap: 1, color: "text.secondary", fontSize: "12px"},
    version: {display: "flex", flexWrap: "wrap", justifyContent: "center", gap: "4px", color: "text.disabled", fontSize: "12px"},
}

type Props = {
    label: string,
    tag: string,
}

export function Footer(props: Props) {
    const {label, tag} = props
    return (
        <Box sx={SX.box}>
            <Box sx={SX.links}>
                {Object.values(IvoryLinks).map(({name, link}, index) => (
                    <Fragment key={name}>
                        {index > 0 && <Box>·</Box>}
                        <Link href={link} target={"_blank"} rel={"noopener"} underline={"hover"} color={"inherit"}>
                            {name}
                        </Link>
                    </Fragment>
                ))}
            </Box>
            <Box sx={SX.version}>
                <Link href={`https://github.com/veegres/ivory/tree/${tag}`} target={"_blank"} rel={"noopener"} underline={"none"} color={"inherit"}>
                    {label}
                </Link>
                <Box>·</Box>
                <Box>
                    © {new Date().getFullYear()} <Link href={"https://veegres.com"} target={"_blank"} rel={"noopener"} underline={"none"} color={"inherit"}>Veegres Labs</Link>
                </Box>
            </Box>
        </Box>
    )
}
