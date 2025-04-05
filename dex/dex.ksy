# https://ide.kaitai.io/#
meta:
  id: dex
  file-extension: dex
  endian: le
seq:
  - id: header
    type: dex_header
instances:
  string_ids:
    io: _root._io
    pos: _root.header.string_ids_offset
    type: u4
    repeat: expr
    repeat-expr: _root.header.string_ids_size / 4
  type_ids:
    io: _root._io
    pos: _root.header.type_ids_offset
    type: u4
    repeat: expr
    repeat-expr: _root.header.type_ids_size / 4
  proto_ids:
    io: _root._io
    pos: _root.header.proto_ids_offset
    type: proto_id_item
    repeat: expr
    repeat-expr: _root.header.proto_ids_size / 12
  field_ids:
    io: _root._io
    pos: _root.header.field_ids_offset
    type: field_id_item
    repeat: expr
    repeat-expr: _root.header.field_ids_size / 12
  method_ids:
    io: _root._io
    pos: _root.header.method_ids_offset
    type: method_id_item
    repeat: expr
    repeat-expr: _root.header.method_ids_size / 12
  class_def_ids:
    io: _root._io
    pos: _root.header.class_defs_offset
    type: class_def_item
    repeat: expr
    repeat-expr: _root.header.class_defs_size / (8*4)
types:
  dex_header:
    seq:
      - id: magic
        contents: [100,101,120,10,48,51,53,0]
      - id: checksum
        type: u4
      - id: signature
        size: 20
      - id: file_size
        type: u4
      - id: header_size
        type: u4
      - id: endiantag
        type: u4
      - id: link_size
        type: u4
      - id: link_offset
        type: u4
      - id: map_offset
        type: u4
      - id: string_ids_size
        type: u4
      - id: string_ids_offset
        type: u4
      - id: type_ids_size
        type: u4
      - id: type_ids_offset
        type: u4
      - id: proto_ids_size
        type: u4
      - id: proto_ids_offset
        type: u4
      - id: field_ids_size
        type: u4
      - id: field_ids_offset
        type: u4
      - id: method_ids_size
        type: u4
      - id: method_ids_offset
        type: u4
      - id: class_defs_size
        type: u4
      - id: class_defs_offset
        type: u4
  proto_id_item:
    seq:
      - id: shorty_idx
        type: u4
      - id: return_type_idx
        type: u4
      - id: parameters_offset
        type: u4
  field_id_item:
    seq:
      - id: class_idx
        type: u4
      - id: type_idx
        type: u4
      - id: name_idx
        type: u4
  method_id_item:
    seq:
      - id: class_idx
        type: u4
      - id: proto_idx
        type: u4
      - id: name_idx
        type: u4
  class_def_item:
    seq:
      - id: class_idx
        type: u4
      - id: access_flags
        type: u4
      - id: super_class_idx
        type: u4
      - id: interfaces_offset
        type: u4
      - id: source_file_idx
        type: u4
      - id: annotations_offset
        type: u4
      - id: class_data_offset
        type: u4
      - id: static_values_offset
        type: u4
