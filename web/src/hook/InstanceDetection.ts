import {InstanceDetection, InstanceMap, Sidecar} from "../app/types";
import {useQuery} from "@tanstack/react-query";
import {instanceApi} from "../app/api";
import {useEffect, useMemo, useRef} from "react";
import {
    combineInstances,
    createInstanceColors,
    getDomain,
} from "../app/utils";

export function useInstanceDetection(name: string, instances: Sidecar[]): InstanceDetection {
    const defaultSidecar = useRef(instances[0])
    const defaultData = useRef<InstanceMap>({})
    const query = useQuery(
        ["instance/overview", name],
        () => instanceApi.overview({...defaultSidecar.current, cluster: name}),
        {retry: false}
    )
    const {errorUpdateCount, refetch, data, isFetching} = query
    defaultData.current = data === undefined ? defaultData.current : data
    const colors = createInstanceColors(defaultData.current)
    const [combinedInstanceMap, warning] = useMemo(() => combineInstances(instances, defaultData.current), [instances])

    useEffect(handleEffectNextRequest, [errorUpdateCount, instances, refetch])

    return {
        defaultInstance: combinedInstanceMap[getDomain(defaultSidecar.current)],
        combinedInstanceMap,
        warning,
        colors,
        fetching: isFetching,
        refetch: handleRefetch,
    }

    /**
     * When we refetch we need to start from scratch, that is why we change
     * `defaultSidecar` value to the first one. After that we clean `useQuery`
     * state and start fetching from scratch with error count equal 0. That
     * is why `handleEffectNextRequest` continue working.
     */
    async function handleRefetch() {
        defaultSidecar.current = instances[0]
        query.remove()
        await query.refetch()
    }

    /**
     * We do refetch in each error, we shouldn't refetch when:
     * - error count is bigger than our amount of instances
     * - when sidecar has default value
     * - when we already fetched this instance, and it still the same
     */
    function handleEffectNextRequest() {
        if (errorUpdateCount < instances.length &&
            defaultSidecar.current !== undefined &&
            instances[errorUpdateCount] !== defaultSidecar.current
        ) {
            defaultSidecar.current = instances[errorUpdateCount]
            refetch().then()
        }
    }
}
