Feature: MusicExpress account
    Users can create and manage accounts

    Scenario Outline: Create and manage account
        Given registration of <s_name> with email <s_email> and password <s_passwd>
        When update profile info to <e_name>, <e_email> and <e_passwd>
        Then get new profile of <e_name> with <e_email> and password <e_passwd>

        Examples:
            | s_name | s_email | s_passwd | e_name | e_email | e_passwd |
            | Pasha  | p@ash.a | pasha    | Egor   | e@go.r  | egor     |
            | John   | j@oh.n  | john     | Alex   | a@le.x  | alex     |
