import {configureStore, ThunkAction, Action} from '@reduxjs/toolkit';
import counterReducer from '../features/w-example/CounterSlice';
import {clusterInfoSlice} from "../features/cluster/ClusterInfoSlice";
import {clusterListSlice} from "../features/cluster/ClusterListSlice";
import {nodePatroniSlice} from "../features/node/NodePatroniSlice";
import {nodeConfigSlice} from "../features/node/NodeConfigSlice";

export const store = configureStore({
    reducer: {
        counter: counterReducer,
        clusterInfo: clusterInfoSlice.reducer,
        clusterList: clusterListSlice.reducer,
        nodePatroni: nodePatroniSlice.reducer,
        nodeConfig: nodeConfigSlice.reducer
    },
});

export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<ReturnType, RootState, unknown, Action<string>>;
