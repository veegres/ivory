import {OpenInNew} from "@mui/icons-material"
import {Box, IconButton} from "@mui/material"

import {List} from "../../../shared/component/box/List"
import {ListItem} from "../../../shared/component/box/ListItem"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {IvoryLinks} from "../../../shared/helper/HelperUtils"

const SX: SxPropsMap = {
    scroll: {display: "flex", flexDirection: "column", padding: "0 15px", gap: 3},
    image: {display: "flex", justifyContent: "center"},
    text: {fontSize: 14, textAlign: "justify"},
}

export function SettingsAbout() {
    return (
        <Box sx={SX.scroll}>
            <Box sx={SX.image}><img src={"/ivory.png"} width={200} height={200} alt={"Ivory"}/></Box>
            <Box sx={SX.text}>
                Ivory is an open-source project designed to simplify and visualize work with Postgres
                clusters. Initially, this tool was developed to ease the life of developers who maintain
                Postgres. But I hope it will help manage and troubleshoot Postgres clusters for both
                developers and database administrators.
            </Box>
            <List name={"Links"}>
                {Object.values(IvoryLinks).map(({name, link}) => (
                    <ListItem key={name} title={name} button={renderLink(link)}/>
                ))}
            </List>
        </Box>
    )

    function renderLink(href: string) {
        return (
            <IconButton href={href} target={"_blank"} rel={"noopener"}>
                <OpenInNew/>
            </IconButton>
        )
    }
}
