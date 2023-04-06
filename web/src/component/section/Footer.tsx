import {Box, Link} from "@mui/material";
import {SxPropsMap} from "../../type/common";
import {IvoryLinks} from "../../app/utils";

const SX: SxPropsMap = {
    box: {
        borderRadius: "15px", fontSize: "12px", fontFamily: "monospace",
        display: "flex", justifyContent: "center", gap: 3, padding: "10px 20px",
    },
}

export function Footer() {
    return (
        <Box sx={SX.box}>
            {Object.values(IvoryLinks).map(({name, link}) => (
                <Link href={link} target={"_blank"} rel={"noopener"} underline={"hover"} color={"inherit"}>
                    {name}
                </Link>
            ))}
        </Box>
    )
}
