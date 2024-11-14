package models

import "gopkg.in/telebot.v3"

// Message
var MessageHomeDefault = `🏠 *Menyu*

🤖 DastyorBot-ga xush kelibsiz! Men matnlar bilan ishlaydigan AI yordamchiman.

Quyidagi xizmatlardan birini tanlang:`

// Buttons
var ButtonHome = &telebot.Btn{Unique: "home", Text: "🏠 Menyu"}
var ButtonRewrite = &telebot.Btn{Unique: "f_Rewrite", Text: "✏️ Qayta yozish"}
var ButtonEnhance = &telebot.Btn{Unique: "f_Enhance", Text: "✏️ Kengaytirish"}
