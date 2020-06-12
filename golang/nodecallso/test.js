var ref = require('ref')
var ffi = require('ffi')
var Struct = require('ref-struct')

var myobj = ref.types.CString
var myobjPtr = ref.refType(myobj)
var GoString = Struct({
  'myobjPtr': myobjPtr,
  'n': 'long'
})

var libm = ffi.Library('hello', {
  'hello': ['string', [GoString]]
})
var name  = 'world'
var gs = new GoString()
gs.n = name.length
gs.myobjPtr = Buffer.alloc(10)
ref.writeCString(gs.myobjPtr, 0, name)

console.log(libm.hello(gs))
