package jotform

import(
    "fmt"
    "strconv"
    "net/http"
    "net/url"   
    "io/ioutil"
    "os"
    "encoding/json"
    "encoding/xml"
    "strings"
    "bytes"
)

const baseURL = "http://api.jotform.com"
const apiVersion = "v1"

type jotformAPIClient struct{
    apiKey string
    outputType string
}

func NewJotFormAPIClient(apiKey string, outputType string) *jotformAPIClient {
    client := &jotformAPIClient{apiKey, strings.ToLower(outputType)}

    return client
}

func (client jotformAPIClient) executeHttpRequest(requestPath string, params interface{}, method string) []byte {

    if client.outputType != "json" {
        requestPath = requestPath + ".xml"
    }

    var path = baseURL + "/" + apiVersion + "/" + requestPath

    var response *http.Response
    var request *http.Request
    var err error

    if method == "GET" {
        path = path + "?" + params.(string)
        request, err = http.NewRequest("GET", path, nil)
        request.Header.Add("apiKey", client.apiKey)
        response, err = http.DefaultClient.Do(request)
    } else if method == "POST" {
        data := params.(map[string]string)
        values := make(url.Values)

        for k, _ := range data {
            values.Set(k, data[k])
        }

        request, err = http.NewRequest("POST", path, strings.NewReader(values.Encode()))
        request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
        request.Header.Add("apiKey", client.apiKey)
        response, err = http.DefaultClient.Do(request)
    } else  if method == "DELETE" {
        request, err = http.NewRequest("DELETE", path, nil)
        request.Header.Add("apiKey", client.apiKey)
        response, err = http.DefaultClient.Do(request)
    } else if method == "PUT" {
        parameters := params.([]byte)

        if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
        } else {
            request, err = http.NewRequest("PUT", path, bytes.NewBuffer(parameters))
            request.Header.Add("apiKey", client.apiKey)
            response, err = http.DefaultClient.Do(request)
        }
    }

    if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
        }

        if client.outputType == "json" {
            var f interface{}
            json.Unmarshal(contents, &f)
            result := f.(map[string]interface{})["content"]
            content, err := json.Marshal(result)

            if err != nil {
                fmt.Printf("%s", err)
                os.Exit(1)
            } else {
                return content   
            }
        } else if client.outputType == "xml" {
            var f interface{}
            xml.Unmarshal(contents, &f)
            return contents
        }
    }

    return nil
}

func createConditions (offset string, limit string, filter map[string]string, orderBy string) string {

    args := make(map[string]interface{})

    args["offset"] = offset
    args["limit"] = limit
    args["filter"] = filter
    args["order_by"] = orderBy

    var params = ""

    for k, _ := range args {
        if k == "filter" {
            if args[k] != nil {
                var value = "{"
                var count = 0

                for key, _ := range filter {
                    count++
                    value = value + "\"" + key + "\":\"" + filter[key] + "\""
                    if count < len(filter) {
                        value = value + ","
                    }
                }
                params = params + "filter=" + value + "}&"
            }
        } else {
            if args[k] != "" {
                params = params + k + "=" + args[k].(string) + "&"
            }
        }
    }
    return params
}

func createHistoryQuery (action string, date string, sortBy string, startDate string, endDate string) string {
    var params = ""

    args := map[string]string {
        "action": action,
        "date": date,
        "sortBy": sortBy,
        "startDate": startDate,
        "endDate": endDate,
    }

    for k, _ := range args {
        if args[k] != "" {
            params = params + k + "=" + args[k] + "&"
        }
    }
    return params
}

//GetUser
//Get user account details for a JotForm user.
//Returns user account type, avatar URL, name, email, website URL and account limits.
func (client jotformAPIClient) GetUser() []byte {
    return client.executeHttpRequest("user", "", "GET")
}

//GetUsage
//Get number of form submissions received this month
//Returns number of submissions, number of SSL form submissions, payment form submissions and upload space used by user.
func (client jotformAPIClient) GetUsage() []byte {
    return client.executeHttpRequest("user/usage", "", "GET")
}

//GetForms
//Get a list of forms for this account
//offset (string): Start of each result set for form list.
//limit (string): Number of results in each result set for form list.
//filter (map[string]string): Filters the query results to fetch a specific form range.
//orderBy (string): Order results by a form field name.
//Returns basic details such as title of the form, when it was created, number of new and total submissions.
func (client jotformAPIClient) GetForms(offset string, limit string, filter map[string]string, orderBy string) []byte {
    var params = createConditions(offset, limit, filter, orderBy)

    return client.executeHttpRequest("user/forms", params, "GET")
}

//GetSubmissions
//Get a list of submissions for this account
//offset (string): Start of each result set for form list.
//limit (string): Number of results in each result set for form list.
//filter (map[string]string): Filters the query results to fetch a specific form range.
//orderBy (string): Order results by a form field name.
//Returns basic details such as title of the form, when it was created, number of new and total submissions.
func (client jotformAPIClient) GetSubmissions(offset string, limit string, filter map[string]string, orderBy string) []byte {
    var params = createConditions(offset, limit, filter, orderBy)

    return client.executeHttpRequest("user/submissions", params, "GET")
}

//GetSubusers
//Get a list of sub users for this account
//Returns list of forms and form folders with access privileges.
func (client jotformAPIClient) GetSubusers() []byte {
    return client.executeHttpRequest("user/subusers", "", "GET")
}

//GetFolders
//Get a list of form folders for this account
//Returns name of the folder and owner of the folder for shared folders.
func (client jotformAPIClient) GetFolders() []byte {
    return client.executeHttpRequest("user/folders", "", "GET")
}

//GetReports
//List of URLS for reports in this account
//Returns reports for all of the forms. ie. Excel, CSV, printable charts, embeddable HTML tables.
func (client jotformAPIClient) GetReports() []byte {
    return client.executeHttpRequest("user/reports", "", "GET")
}

//GetSettings
//Get user's settings for this account
//Returns user's time zone and language.
func (client jotformAPIClient) GetSettings() []byte {
    return client.executeHttpRequest("user/settings", "", "GET")
}

//GetHistory
//Get user activity log
//action (string): Filter results by activity performed. Default is 'all'.
//date (string): Limit results by a date range. If you'd like to limit results by specific dates you can use startDate and endDate fields instead.
//sortBy (string): Lists results by ascending and descending order.
//startDate (string): Limit results to only after a specific date. Format: MM/DD/YYYY.
//endDate (string): Limit results to only before a specific date. Format: MM/DD/YYYY.
//Returns activity log about things like forms created/modified/deleted, account logins and other operations.
func (client jotformAPIClient) GetHistory(action string, date string, sortBy string, startDate string, endDate string) []byte {
    var params = createHistoryQuery(action, date, sortBy, startDate, endDate)

    return client.executeHttpRequest("user/history", params, "GET")
}

//GetForm
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//Returns form ID, status, update and creation dates, submission count etc.
func (client jotformAPIClient) GetForm(formID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10), "", "GET")
}

//GetFormQuestions
//Get a list of all questions on a form.
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//Returns question properties of a form.
func (client jotformAPIClient) GetFormQuestions(formID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/questions", "", "GET")
}

//GetFormQuestion
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//qid (int): Identifier for each question on a form. You can get a list of question IDs from /form/{id}/questions.
//Returns question properties like required and validation.
func (client jotformAPIClient) GetFormQuestion(formID int64, qid int) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/question/" + strconv.Itoa(qid), "", "GET")
}

//GetFormSubmission
//List of a form submissions.
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//offset (string): Start of each result set for form list.
//limit (string): Number of results in each result set for form list.
//filter (map[string]string): Filters the query results to fetch a specific form range.
//orderBy (string): Order results by a form field name.
//Returns submissions of a specific form.
func (client jotformAPIClient) GetFormSubmissions(formID int64, offset string, limit string, filter map[string]string, orderBy string) []byte {
    var params = createConditions(offset, limit, filter, orderBy)

    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/submissions", params, "GET")
}

//CreateFormSubmissions
//Submit data to this form using the API
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//submission (map[string]string): Submission data with question IDs.
//Returns posted submission ID and URL.
func (client jotformAPIClient) CreateFormSubmission(formId int64, submission map[string]string) []byte {
    data := make(map[string]string)

    for k, _ := range submission {
        if strings.Contains(k, "_") {
            data["submission[" + k[0:strings.Index(k, "_")] + "][" + k[strings.Index(k, "_")+1:len(k)] + "]"] = submission[k]
        } else {
            data["submission[" + k + "]"] = submission[k]   
        }
    }

    return client.executeHttpRequest("form/" + strconv.FormatInt(formId, 10) + "/submissions", data, "POST")
}

//GetFormFiles
//List of files uploaded on a form
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//Returns uploaded file information and URLs on a specific form.
func (client jotformAPIClient) GetFormFiles(formID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/files", "", "GET")
}

//GetFormWebhooks
//Get list of webhooks for a form
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//Returns list of webhooks for a specific form.
func (client jotformAPIClient) GetFormWebhooks(formID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/webhooks", "", "GET")
}

//CreateFormWebhook
//Add a new webhook
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//webhookURL (string): Webhook URL is where form data will be posted when form is submitted.
//Returns list of webhooks for a specific form.
func (client jotformAPIClient) CreateFormWebhook(formId int64, webhookURL string) []byte {
    params := map[string]string {
        "webhookURL": webhookURL,
    }

    return client.executeHttpRequest("form/" + strconv.FormatInt(formId, 10) + "/webhooks", params, "POST")
}

//GetSubmission
//Get submission data
//sid (int64): You can get submission IDs when you call /form/{id}/submissions.
//Returns information and answers of a specific submission.
func(client jotformAPIClient) GetSubmission(sid int64) []byte {
    return client.executeHttpRequest("user/submission/" + strconv.FormatInt(sid, 10), "","GET")
}

//GetReport
//Get report details
//reportID (int64): You can get a list of reports from /user/reports.
//Returns properties of a speceific report like fields and status.
func(client jotformAPIClient) GetReport(reportID int64) []byte {
    return client.executeHttpRequest("user/report/" + strconv.FormatInt(reportID, 10), "", "GET")
}

//GetFolder
//folderID (int64): You can get a list of folders from /user/folders.
//Returns a list of forms in a folder, and other details about the form such as folder color.
func (client jotformAPIClient) GetFolder(folderID int64) []byte {
    return client.executeHttpRequest("user/folder/" + strconv.FormatInt(folderID, 10), "", "GET")
}

//GetFormProperties
//Get a list of all properties on a form
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//Returns form properties like width, expiration date, style etc.
func (client jotformAPIClient) GetFormProperties(formID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/properties", "", "GET")
}

//GetFormProperty
//Get a specific property of the form.]
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//propertyKey (string): You can get property keys when you call /form/{id}/properties.
//Returns given property key value.
func (client jotformAPIClient) GetFormProperty(formID int64, propertyKey string) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/properties/" + propertyKey, "", "POST")
}

//DeleteSubmission
//Delete a single submission
//sid (int64): You can get submission IDs when you call /form/{id}/submissions.
//Returns status of request.
func (client jotformAPIClient) DeleteSubmission(sid int64) []byte {
    return client.executeHttpRequest("submission/" + strconv.FormatInt(sid, 10), nil, "DELETE")
}

//EditSubmission
//Edit a single submission
//sid (int64): You can get submission IDs when you call /form/{id}/submissions.
//submission (map[string]string): New submission data with question IDs.
//Returns status of request.
func (client jotformAPIClient) EditSubmission(sid int64, submission map[string]string) []byte {
    data := make(map[string]string)

    for k, _ := range submission {
        if strings.Contains(k, "_") {
            data["submission[" + k[0:strings.Index(k, "_")] + "][" + k[strings.Index(k, "_")+1:len(k)] + "]"] = submission[k]
        } else {
            data["submission[" + k + "]"] = submission[k]   
        }
    }

    return client.executeHttpRequest("submission/" + strconv.FormatInt(sid, 10), data, "POST")
}

//CloneForm
//Clone a single form.
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//Returns status of request.
func (client jotformAPIClient) CloneForm(formID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/clone", nil, "POST")
}

//DeleteFormQuestion
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//qid (int): Identifier for each question on a form. You can get a list of question IDs from /form/{id}/questions.
//Returns status of request.
func (client jotformAPIClient) DeleteFormQuestion(formID int64, qid int) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/question/" + strconv.Itoa(qid), nil, "DELETE")
}

//CreateFormQuestion
//Add new question to specified form.
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//questionProperties (map[string]string): New question properties like type and text.
//Returns properties of new question.
func (client jotformAPIClient) CreateFormQuestion(formID int64, questionProperties map[string]string) []byte {
    question := make(map[string]string)

    for k, _ := range questionProperties {
        question["question[" + k + "]"] = questionProperties[k]
    }

    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/questions", question, "POST")
}

//CreateFormQuestion
//Add new question to specified form.
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//questions ([]byte): New question properties like type and text.
//Returns properties of new question.
func (client jotformAPIClient) CreateFormQuestions(formID int64, questions []byte) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/questions", questions, "PUT")
}

//EditFormQuestion
//Add or edit a single question properties
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//qid (int): Identifier for each question on a form. You can get a list of question IDs from /form/{id}/questions.
//questionProperties (map[string]string): New question properties like type and text.
//Returns edited property and type of question.
func (client jotformAPIClient) EditFormQuestion(formID int64, qid int, questionProperties map[string]string) []byte {
    question := make(map[string]string)

    for k, _ := range questionProperties {
        question["question[" + k + "]"] = questionProperties[k]
    }

    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/question/" + strconv.Itoa(qid), question, "POST")
}

//SetFormProperties
//Add or edit properties of a specific form
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//formProperties (map[string]string): New properties like label width.
//Returns edited properties.
func (client jotformAPIClient) SetFormProperties(formID int64, formProperties map[string]string) []byte {
    properties := make(map[string]string)

    for k, _ := range formProperties {
        properties["properties[" + k + "]"] = formProperties[k]
    }

    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/properties", properties, "POST")
}

//SetFormProperties
//Add or edit properties of a specific form
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//formProperties ([]byte): New properties like label width.
//Returns edited properties.
func (client jotformAPIClient) SetMultipleFormProperties(formID int64, formProperties []byte) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/properties", formProperties, "PUT")
}

//CreateForm
//Create a new form
//form ([]byte): Questions, properties and emails of new form.
//Returns new form.
func (client jotformAPIClient) CreateForms(form []byte) []byte {
    return client.executeHttpRequest("user/forms", form, "PUT")
}

//DeleteForm
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//Returns properties of deleted form.
func (client jotformAPIClient) DeleteForm(formID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10), nil, "DELETE")
}

