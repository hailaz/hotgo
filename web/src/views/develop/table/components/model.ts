// 字段数据类型选项（分组展示）
export const dataTypeOptions = [
  {
    type: 'group',
    label: '整数类型',
    key: 'integer',
    children: [
      { label: 'tinyint', value: 'tinyint' },
      { label: 'smallint', value: 'smallint' },
      { label: 'mediumint', value: 'mediumint' },
      { label: 'int', value: 'int' },
      { label: 'bigint', value: 'bigint' },
    ],
  },
  {
    type: 'group',
    label: '字符串类型',
    key: 'string',
    children: [
      { label: 'char', value: 'char' },
      { label: 'varchar', value: 'varchar' },
      { label: 'tinytext', value: 'tinytext' },
      { label: 'text', value: 'text' },
      { label: 'mediumtext', value: 'mediumtext' },
      { label: 'longtext', value: 'longtext' },
    ],
  },
  {
    type: 'group',
    label: '日期时间',
    key: 'datetime',
    children: [
      { label: 'date', value: 'date' },
      { label: 'datetime', value: 'datetime' },
      { label: 'timestamp', value: 'timestamp' },
      { label: 'time', value: 'time' },
    ],
  },
  {
    type: 'group',
    label: '数值类型',
    key: 'numeric',
    children: [
      { label: 'decimal', value: 'decimal' },
      { label: 'float', value: 'float' },
      { label: 'double', value: 'double' },
    ],
  },
  {
    type: 'group',
    label: '其他类型',
    key: 'other',
    children: [
      { label: 'json', value: 'json' },
      { label: 'blob', value: 'blob' },
      { label: 'enum', value: 'enum' },
    ],
  },
];

// 索引类型选项
export const indexTypeOptions = [
  { label: '普通索引 (INDEX)', value: 'INDEX' },
  { label: '唯一索引 (UNIQUE)', value: 'UNIQUE' },
  { label: '全文索引 (FULLTEXT)', value: 'FULLTEXT' },
];

// 存储引擎选项
export const engineOptions = [
  { label: 'InnoDB', value: 'InnoDB' },
  { label: 'MyISAM', value: 'MyISAM' },
];

// 整数类型列表（用于判断是否显示无符号选项）
export const integerTypes = ['tinyint', 'smallint', 'mediumint', 'int', 'bigint'];

// 需要小数位的类型
export const decimalTypes = ['decimal', 'float', 'double'];

// 需要长度的类型
export const lengthTypes = ['char', 'varchar', 'tinyint', 'smallint', 'mediumint', 'int', 'bigint', 'decimal', 'float', 'double'];

// 默认字段
export interface ColumnItem {
  name: string;
  dataType: string;
  length: number;
  decimal: number;
  isNullable: boolean;
  defaultValue: string;
  comment: string;
  isPrimaryKey: boolean;
  isAutoInc: boolean;
  isUnsigned: boolean;
}

// 默认索引
export interface IndexItem {
  name: string;
  type: string;
  columns: string[];
}

// 新建空字段
export function newColumn(): ColumnItem {
  return {
    name: '',
    dataType: 'varchar',
    length: 255,
    decimal: 0,
    isNullable: true,
    defaultValue: '',
    comment: '',
    isPrimaryKey: false,
    isAutoInc: false,
    isUnsigned: false,
  };
}

// 新建空索引
export function newIndex(): IndexItem {
  return {
    name: '',
    type: 'INDEX',
    columns: [],
  };
}

// 预设字段模板（基于系统全表字段统计，按使用频率排序）
export const presetColumns: Record<string, ColumnItem[]> = {
  id主键: [
    {
      name: 'id',
      dataType: 'bigint',
      length: 20,
      decimal: 0,
      isNullable: false,
      defaultValue: '',
      comment: '主键ID',
      isPrimaryKey: true,
      isAutoInc: true,
      isUnsigned: true,
    },
  ],
  时间戳字段: [
    {
      name: 'created_at',
      dataType: 'datetime',
      length: 0,
      decimal: 0,
      isNullable: true,
      defaultValue: '',
      comment: '创建时间',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
    {
      name: 'updated_at',
      dataType: 'datetime',
      length: 0,
      decimal: 0,
      isNullable: true,
      defaultValue: '',
      comment: '修改时间',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  状态字段: [
    {
      name: 'status',
      dataType: 'tinyint',
      length: 1,
      decimal: 0,
      isNullable: false,
      defaultValue: '1',
      comment: '状态',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  排序字段: [
    {
      name: 'sort',
      dataType: 'int',
      length: 11,
      decimal: 0,
      isNullable: false,
      defaultValue: '0',
      comment: '排序',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  删除标记: [
    {
      name: 'deleted_at',
      dataType: 'datetime',
      length: 0,
      decimal: 0,
      isNullable: true,
      defaultValue: '',
      comment: '删除时间',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  备注字段: [
    {
      name: 'remark',
      dataType: 'varchar',
      length: 255,
      decimal: 0,
      isNullable: true,
      defaultValue: '',
      comment: '备注',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  名称字段: [
    {
      name: 'name',
      dataType: 'varchar',
      length: 100,
      decimal: 0,
      isNullable: false,
      defaultValue: '',
      comment: '名称',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  标题字段: [
    {
      name: 'title',
      dataType: 'varchar',
      length: 255,
      decimal: 0,
      isNullable: false,
      defaultValue: '',
      comment: '标题',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  描述字段: [
    {
      name: 'description',
      dataType: 'varchar',
      length: 255,
      decimal: 0,
      isNullable: true,
      defaultValue: '',
      comment: '描述',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  树形关系字段: [
    {
      name: 'pid',
      dataType: 'bigint',
      length: 20,
      decimal: 0,
      isNullable: false,
      defaultValue: '0',
      comment: '上级ID',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
    {
      name: 'level',
      dataType: 'int',
      length: 11,
      decimal: 0,
      isNullable: false,
      defaultValue: '1',
      comment: '关系树等级',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
    {
      name: 'tree',
      dataType: 'varchar',
      length: 512,
      decimal: 0,
      isNullable: true,
      defaultValue: '',
      comment: '关系树',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  操作人字段: [
    {
      name: 'created_by',
      dataType: 'bigint',
      length: 20,
      decimal: 0,
      isNullable: true,
      defaultValue: '0',
      comment: '创建者',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
    {
      name: 'updated_by',
      dataType: 'bigint',
      length: 20,
      decimal: 0,
      isNullable: true,
      defaultValue: '0',
      comment: '更新者',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  管理员ID: [
    {
      name: 'member_id',
      dataType: 'bigint',
      length: 20,
      decimal: 0,
      isNullable: false,
      defaultValue: '0',
      comment: '管理员ID',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  分类ID: [
    {
      name: 'category_id',
      dataType: 'bigint',
      length: 20,
      decimal: 0,
      isNullable: true,
      defaultValue: '0',
      comment: '分类ID',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  内容字段: [
    {
      name: 'content',
      dataType: 'longtext',
      length: 0,
      decimal: 0,
      isNullable: true,
      defaultValue: '',
      comment: '内容',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  图片字段: [
    {
      name: 'image',
      dataType: 'varchar',
      length: 255,
      decimal: 0,
      isNullable: true,
      defaultValue: '',
      comment: '图片',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
    {
      name: 'images',
      dataType: 'json',
      length: 0,
      decimal: 0,
      isNullable: true,
      defaultValue: '',
      comment: '多图',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  附件字段: [
    {
      name: 'attachfile',
      dataType: 'varchar',
      length: 255,
      decimal: 0,
      isNullable: true,
      defaultValue: '',
      comment: '附件',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
    {
      name: 'attachfiles',
      dataType: 'json',
      length: 0,
      decimal: 0,
      isNullable: true,
      defaultValue: '',
      comment: '多附件',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  金额字段: [
    {
      name: 'money',
      dataType: 'decimal',
      length: 10,
      decimal: 2,
      isNullable: false,
      defaultValue: '0.00',
      comment: '金额',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  价格字段: [
    {
      name: 'price',
      dataType: 'decimal',
      length: 10,
      decimal: 2,
      isNullable: false,
      defaultValue: '0.00',
      comment: '价格',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  联系方式字段: [
    {
      name: 'email',
      dataType: 'varchar',
      length: 60,
      decimal: 0,
      isNullable: true,
      defaultValue: '',
      comment: '邮箱',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
    {
      name: 'mobile',
      dataType: 'varchar',
      length: 20,
      decimal: 0,
      isNullable: true,
      defaultValue: '',
      comment: '手机号码',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  IP地址: [
    {
      name: 'ip',
      dataType: 'varchar',
      length: 128,
      decimal: 0,
      isNullable: true,
      defaultValue: '',
      comment: 'IP地址',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  订单号字段: [
    {
      name: 'order_sn',
      dataType: 'varchar',
      length: 64,
      decimal: 0,
      isNullable: true,
      defaultValue: '',
      comment: '订单号',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  开关字段: [
    {
      name: 'switch',
      dataType: 'tinyint',
      length: 1,
      decimal: 0,
      isNullable: true,
      defaultValue: '0',
      comment: '开关',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  时间范围字段: [
    {
      name: 'start_at',
      dataType: 'datetime',
      length: 0,
      decimal: 0,
      isNullable: true,
      defaultValue: '',
      comment: '开始时间',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
    {
      name: 'end_at',
      dataType: 'datetime',
      length: 0,
      decimal: 0,
      isNullable: true,
      defaultValue: '',
      comment: '结束时间',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  浏览次数: [
    {
      name: 'views',
      dataType: 'bigint',
      length: 20,
      decimal: 0,
      isNullable: true,
      defaultValue: '0',
      comment: '浏览次数',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  性别字段: [
    {
      name: 'sex',
      dataType: 'tinyint',
      length: 1,
      decimal: 0,
      isNullable: true,
      defaultValue: '1',
      comment: '性别',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  头像字段: [
    {
      name: 'avatar',
      dataType: 'varchar',
      length: 255,
      decimal: 0,
      isNullable: true,
      defaultValue: '',
      comment: '头像',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
  动态键值对: [
    {
      name: 'map',
      dataType: 'json',
      length: 0,
      decimal: 0,
      isNullable: true,
      defaultValue: '',
      comment: '动态键值对',
      isPrimaryKey: false,
      isAutoInc: false,
      isUnsigned: false,
    },
  ],
};

// 预设字段描述信息
export interface PresetMeta {
  key: string;
  label: string;
  description: string;
  fields: string;
  tag: string;
}

export const presetMeta: PresetMeta[] = [
  { key: 'id主键', label: 'ID 主键', description: '自增主键，几乎所有表必备', fields: 'id (bigint)', tag: '必备' },
  { key: '时间戳字段', label: '时间戳', description: '记录创建和修改时间，出现频率最高', fields: 'created_at, updated_at', tag: '必备' },
  { key: '删除标记', label: '软删除', description: '软删除标记，GoFrame 约定字段', fields: 'deleted_at (datetime)', tag: '必备' },
  { key: '状态字段', label: '状态', description: '通用状态标识，1=正常 2=禁用等', fields: 'status (tinyint)', tag: '常用' },
  { key: '排序字段', label: '排序', description: '列表排序权重，值越小越靠前', fields: 'sort (int)', tag: '常用' },
  { key: '备注字段', label: '备注', description: '通用备注信息，出现在17张表中', fields: 'remark (varchar 255)', tag: '常用' },
  { key: '名称字段', label: '名称', description: '通用名称，出现在13张表中', fields: 'name (varchar 100)', tag: '常用' },
  { key: '标题字段', label: '标题', description: '内容标题，适用于文章/公告等场景', fields: 'title (varchar 255)', tag: '常用' },
  { key: '描述字段', label: '描述', description: '简短描述信息', fields: 'description (varchar 255)', tag: '常用' },
  { key: '操作人字段', label: '操作人', description: '记录创建者和更新者的用户ID', fields: 'created_by, updated_by (bigint)', tag: '审计' },
  { key: '管理员ID', label: '管理员ID', description: '关联后台管理员', fields: 'member_id (bigint)', tag: '关联' },
  { key: '分类ID', label: '分类ID', description: '关联分类表', fields: 'category_id (bigint)', tag: '关联' },
  { key: '树形关系字段', label: '树形关系', description: '父子层级关系，适用于菜单/分类等', fields: 'pid, level, tree', tag: '结构' },
  { key: '内容字段', label: '富文本内容', description: '长文本内容，如文章正文、详情等', fields: 'content (longtext)', tag: '内容' },
  { key: '图片字段', label: '图片', description: '单图 + 多图，支持JSON存储', fields: 'image (varchar), images (json)', tag: '媒体' },
  { key: '附件字段', label: '附件', description: '单附件 + 多附件', fields: 'attachfile (varchar), attachfiles (json)', tag: '媒体' },
  { key: '头像字段', label: '头像', description: '用户/管理员头像', fields: 'avatar (varchar 255)', tag: '用户' },
  { key: '性别字段', label: '性别', description: '1=男 2=女 等', fields: 'sex (tinyint)', tag: '用户' },
  { key: '联系方式字段', label: '联系方式', description: '邮箱和手机号', fields: 'email (varchar), mobile (varchar)', tag: '用户' },
  { key: '金额字段', label: '金额', description: '高精度金额，两位小数', fields: 'money (decimal 10,2)', tag: '业务' },
  { key: '价格字段', label: '价格', description: '商品/服务价格', fields: 'price (decimal 10,2)', tag: '业务' },
  { key: '订单号字段', label: '订单号', description: '业务订单编号', fields: 'order_sn (varchar 64)', tag: '业务' },
  { key: '开关字段', label: '开关', description: '布尔开关 0/1', fields: 'switch (tinyint)', tag: '业务' },
  { key: '时间范围字段', label: '时间范围', description: '开始-结束时间区间', fields: 'start_at, end_at (datetime)', tag: '业务' },
  { key: 'IP地址', label: 'IP地址', description: '记录客户端IP', fields: 'ip (varchar 128)', tag: '日志' },
  { key: '浏览次数', label: '浏览次数', description: '内容浏览/点击计数', fields: 'views (bigint)', tag: '统计' },
  { key: '动态键值对', label: '动态键值对', description: 'JSON格式的自由扩展字段', fields: 'map (json)', tag: '扩展' },
];

// 预设分组（用于下拉菜单分组展示）
export const presetGroups = [
  { label: '基础必备', tags: ['必备'] },
  { label: '常用字段', tags: ['常用'] },
  { label: '用户相关', tags: ['用户'] },
  { label: '关联 & 结构', tags: ['关联', '结构', '审计'] },
  { label: '内容 & 媒体', tags: ['内容', '媒体'] },
  { label: '业务字段', tags: ['业务'] },
  { label: '其他', tags: ['日志', '统计', '扩展'] },
];

// 兼容旧的 presetOptions（保留导出以防其他地方使用）
export const presetOptions = presetMeta.map((m) => ({
  label: m.label,
  value: m.key,
}));

// 字段排序权重表：添加预设后按此权重排序，保证字段顺序规整
// 权重越小越靠前；未在表中的自定义字段默认 500（居中）
export const fieldSortWeight: Record<string, number> = {
  // === 主键，永远第一 ===
  id: 0,

  // === 树形/层级关系，紧跟主键 ===
  pid: 10,
  level: 11,
  tree: 12,

  // === 关联ID ===
  member_id: 50,
  category_id: 55,

  // === 核心业务字段（靠前） ===
  name: 100,
  title: 110,
  description: 120,
  content: 130,

  // === 联系方式 ===
  email: 200,
  mobile: 210,
  avatar: 220,
  sex: 230,
  birthday: 240,

  // === 媒体/附件 ===
  image: 300,
  images: 301,
  attachfile: 310,
  attachfiles: 311,

  // === 金额/价格/订单 ===
  money: 400,
  price: 410,
  order_sn: 420,

  // === 其他业务 ===
  ip: 500,
  views: 510,
  switch: 520,
  start_at: 530,
  end_at: 531,
  map: 550,

  // === 通用状态/排序/备注（靠后） ===
  status: 800,
  sort: 810,
  remark: 820,

  // === 审计字段 ===
  created_by: 850,
  updated_by: 851,

  // === 时间戳（最后） ===
  created_at: 900,
  updated_at: 910,
  deleted_at: 950,
};

// 获取字段排序权重（未定义的字段返回默认值 500）
export function getFieldWeight(fieldName: string): number {
  return fieldSortWeight[fieldName] ?? 500;
}
