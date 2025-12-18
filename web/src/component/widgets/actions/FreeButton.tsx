import {useRouterFree} from "../../../api/management/hook"
import {AlertButton} from "../../view/button/AlertButton"

export function FreeButton() {
    const free = useRouterFree()

    return (
        <AlertButton
            variant={"outlined"}
            size={"small"}
            label={"Free"}
            title={"Free space"}
            description={`Ivory can take up a significant amount of disk space over time. This
             will clean up query logs and bloat logs.`}
            loading={free.isPending}
            onClick={() => free.mutate()}
        />
    )
}
