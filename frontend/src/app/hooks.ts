import { TypedUseSelectorHook, useDispatch, useSelector } from 'react-redux';
import type { RootState, AppDispatch } from './store';

// Use throughout your app instead of plain `useDispatch` and `useSelector`
export const useAppDispatch = () => useDispatch<AppDispatch>();
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;
export const restClient = <T>(url: string): Promise<T> => fetch(`/api${url}`).then(response => {
    if (!response.ok) throw new Error(response.statusText)
    return response.json()
})
