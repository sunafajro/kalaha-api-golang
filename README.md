#API для ресурса кала-ха.рф
##Учебные пособия:
 - GET /api/books - схема всех учебных пособий.
 - GET /api/books/:id - схема указанного учебного пособия. id - номер пособия.
 - GET /api/books/:id/:chapter - схема главы учебного пособия. id - номер пособия, chapter - номер главы.
 - GET /api/books/:id/file/:name/:type - файл указанного учебного пособия. id - номер пособия, name - имя файла, type - тип файла