import {AlertButton} from "../../../shared/component/button/AlertButton"
import {useRouterFree} from "../api/ManagementHook"

export function ManageFreeButton() {
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
