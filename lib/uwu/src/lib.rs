use std::ffi::{CStr, CString};
use std::os::raw;

#[no_mangle]
pub extern "C" fn uwuify(s: *const raw::c_char) -> *const raw::c_char {
    let buf = unsafe { CStr::from_ptr(s).to_bytes() };
    let mut temp1 = vec![0u8; uwuifier::round_up16(buf.len()) * 16];
    let mut temp2 = vec![0u8; uwuifier::round_up16(buf.len()) * 16];
    let uwu = uwuifier::uwuify_sse(buf, &mut temp1, &mut temp2);
    let uwu = CString::new(uwu).unwrap(); // won't have \0 here
    uwu.into_raw()
}
