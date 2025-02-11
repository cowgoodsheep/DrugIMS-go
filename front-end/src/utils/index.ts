import { setMax as setMaxAction, setOpen as setOpenAction, setType as setTypeAction } from '../store/reducers/model';
import { useAppDispatch, useAppSelector } from './declare'
export const useModel = () => {
  const model = useAppSelector(state => state.model)
  const dispath = useAppDispatch()
  const { type, open, titleMap } = model
  const setType = (type: number) => {
    dispath(setTypeAction(type))
    setOpen(true)
  }
  const setOpen = (state: boolean) => {
    dispath(setOpenAction(state))
  }
  return {
    setOpen, setType, type, open, titleMap
  }
}
export const addOneDay = (dateStr) => {
  var parts = dateStr.split('-');
  var year = parseInt(parts[0], 10);
  var month = parseInt(parts[1], 10) - 1;  // 月份从0开始计数
  var day = parseInt(parts[2], 10);

  var date = new Date(year, month, day);
  date.setDate(date.getDate() + 1);

  var nextDay = date.getFullYear() + '-'
    + ('0' + (date.getMonth() + 1)).slice(-2) + '-'
    + ('0' + date.getDate()).slice(-2);

  return nextDay;
}

export const formatDate = (isoString) => {
  // 创建一个 Date 对象
  const date = new Date(isoString);

  // 获取年、月、日
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0'); // 月份从0开始，需要加1
  const day = String(date.getDate()).padStart(2, '0');

  // 获取小时、分钟、秒
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');

  // 拼接成目标格式
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
};