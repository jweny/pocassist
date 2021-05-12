const List = [
  {
    id: 1,
    vul_name: "loudong1",
    xml: "xxxx1",
    type: "SQL注入",
    module: "mokuai1",
    level: "high"
  },
  {
    id: 2,
    vul_name: "loudong1",
    xml: "xxxx1",
    type: "SQL注入",
    module: "mokuai1",
    level: "high"
  },
  {
    id: 3,
    vul_name: "loudong1",
    xml: "xxxx1",
    type: "SQL注入",
    module: "mokuai1",
    level: "high"
  },
  {
    id: 4,
    vul_name: "loudong1",
    xml: "xxxx1",
    type: "SQL注入",
    module: "mokuai1",
    level: "high"
  },
  {
    id: 5,
    vul_name: "loudong1",
    xml: "xxxx1",
    type: "SQL注入",
    module: "mokuai1",
    level: "high"
  },
  {
    id: 6,
    vul_name: "loudong1",
    xml: "xxxx1",
    type: "SQL注入",
    module: "mokuai1",
    level: "high"
  },
  {
    id: 7,
    vul_name: "loudong1",
    xml: "xxxx1",
    type: "SQL注入",
    module: "mokuai1",
    level: "high"
  },
  {
    id: 8,
    vul_name: "loudong1",
    xml: "xxxx1",
    type: "SQL注入",
    module: "mokuai1",
    level: "high"
  },
  {
    id: 9,
    vul_name: "loudong1",
    xml: "xxxx1",
    type: "SQL注入",
    module: "mokuai1",
    level: "high"
  },
  {
    id: 10,
    vul_name: "loudong1",
    xml: "xxxx1",
    type: "SQL注入",
    module: "mokuai1",
    level: "high"
  },
  {
    id: 11,
    vul_name: "loudong1",
    xml: "xxxx1",
    type: "SQL注入",
    module: "mokuai1",
    level: "high"
  },
  {
    id: 12,
    vul_name: "loudong1",
    xml: "xxxx1",
    type: "SQL注入",
    module: "mokuai1",
    level: "high"
  },
  {
    id: 13,
    vul_name: "loudong1",
    xml: "xxxx1",
    type: "SQL注入",
    module: "mokuai1",
    level: "high"
  },
  {
    id: 14,
    vul_name: "loudong1",
    xml: "xxxx1",
    type: "SQL注入",
    module: "mokuai1",
    level: "high"
  },
  {
    id: 15,
    vul_name: "loudong1",
    xml: "xxxx1",
    type: "SQL注入",
    module: "mokuai1",
    level: "high"
  }
];
export default List;

const result = {
  tests: {
    test: {
      request: {
        method: "GET",
        url: "$(scheme)://$(host):$(port)$(path)",
        version: "HTTP/1.0",
        post_text: null,
        cookies: null,
        custom_headers: null
      },
      response: {
        var: {
          "@description": "",
          "@name": "response_code",
          "@source": "statusline",
          "#text": "^.*\\s(\\d\\d\\d)\\s"
        }
      },
      criteria: {
        "@operator": "OR",
        criterion: [
          {
            "@comment": "comment",
            "@operator": "pattern match",
            "@value": '(?is)<a\\s+?href="[^"]+?\\.action[^"]*?"[^>]+?>',
            "@variable": "$(body)"
          },
          {
            "@comment": "comment",
            "@operator": "pattern match",
            "@value":
              '(?is)<form\\s*?[^>]+?\\s*?action="[^"]+?\\.action"\\s+?method="post">',
            "@variable": "$(body)"
          }
        ]
      }
    }
  }
};
