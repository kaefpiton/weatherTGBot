package domain

// todo исправить названия
const Err_usr_state = "error"

// users group
const User_unauth_state = "unauthorised"
const User_auth_state = "authorised"

// todo возможно выпилить
// const User_choice_state = "user_choice"
const User_city_choice_state = "user_city_choice"

// admin group
// todo добавить мапу с админами (юзернеймами)*
const Admin_auth_state = "admin_auth"
const Admin_enter_state = "admin_enter"
const Admin_set_sticker_state = "admin_set_sticker"

// todo админ выбирает категорию стикера
const Admin_set_sticker_category = "admin_set_sticker_category"
