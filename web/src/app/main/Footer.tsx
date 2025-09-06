import {Box, Link} from "@mui/material";
import {IvoryLinks} from "../utils";
import {SxPropsMap} from "../type";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", fontFamily: "monospace", margin: "5px 20px 5px", gap: "5px"},
    links: {display: "flex", justifyContent: "center", gap: 3, fontSize: "12px"},
    version: {display: "flex", justifyContent: "center", color: "text.disabled", fontSize: "8px"},
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
                {Object.values(IvoryLinks).map(({name, link}) => (
                    <Link key={name} href={link} target={"_blank"} rel={"noopener"} underline={"hover"} color={"inherit"}>
                        {name}
                    </Link>
                ))}
            </Box>
            <Box sx={SX.version}>
                <Link href={`https://github.com/veegres/ivory/tree/${tag}`} target={"_blank"} rel={"noopener"} underline={"none"} color={"inherit"}>
                    {label}
                </Link>
            </Box>
        </Box>
    )
}
