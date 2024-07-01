import {Box, Link} from "@mui/material";
import {SxPropsMap} from "../../type/general";
import {IvoryLinks, shortUuid} from "../../app/utils";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", fontFamily: "monospace", margin: "20px 20px 5px", gap: "5px"},
    links: {display: "flex", justifyContent: "center", gap: 3, fontSize: "12px"},
    version: {display: "flex", justifyContent: "center", color: "text.disabled", fontSize: "8px"},
}

type Props = {
    tag: string,
    commit: string,
}

export function Footer(props: Props) {
    const {tag, commit} = props
    return (
        <Box sx={SX.box}>
            <Box sx={SX.links}>
                {Object.values(IvoryLinks).map(({name, link}) => (
                    <Link key={name} href={link} target={"_blank"} rel={"noopener"} underline={"hover"} color={"inherit"}>
                        {name}
                    </Link>
                ))}
            </Box>
            <Box sx={SX.version}>{`Ivory: ${tag} (${commit})`}</Box>
        </Box>
    )
}
