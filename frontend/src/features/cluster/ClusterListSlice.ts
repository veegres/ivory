import {createAsyncThunk, createSlice} from '@reduxjs/toolkit';
import {RootState} from '../../app/store';
import {restClient} from "../../app/hooks";

export interface ClusterList {
    name: string,
    nodes: string[]
}

export interface ClusterListState {
    data: ClusterList[];
    status: 'idle' | 'loading' | 'failed';
}

const initialState: ClusterListState = {
    data: [],
    status: 'idle',
};

export const incrementAsync = createAsyncThunk(
    'cluster/fetchList',
    async () => await restClient<ClusterList[]>(`/cluster/list`)
);

export const clusterListSlice = createSlice({
    name: 'cluster/list',
    initialState,
    reducers: {},
    extraReducers: (builder) => {
        builder
            .addCase(incrementAsync.pending, (state) => {
                state.status = 'loading';
            })
            .addCase(incrementAsync.fulfilled, (state, action) => {
                state.status = 'idle';
                state.data = action.payload;
            });
    },
});

export const selectClusterListData = (state: RootState) => state.clusterList.data;
