namespace go example.store

struct MultiTypeReq {
  1: required i64 para1,
  2: i32 para2,
  3: string para3
  4: required bool para4
}

struct Resp {
  1: required i64 para1,
  2: string para3
  3: bool para4
}

 service HertzService {
     Resp Method1(1: MultiTypeReq request) ( // 处理请求参数需要对请求参数进行遍历，遍历他的每一个 field，从而映射为不同参数
         api.get="/company/department/group/user:id/name";
         api.consume="json";
      );
 } (
    api.base_domain = "127.0.0.1:8888";
    api.version = "1.0.0";
    api.scheme = "https";
 )