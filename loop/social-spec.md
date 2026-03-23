# Social Layer Specification

**Using Code Graph primitives to formally describe four app modes that together encompass all major platform communication patterns.**

Matt Searles + Claude · March 2026

---

## Architecture

The Social layer is not one app. It's four **modes** — each a distinct composition of Code Graph primitives, each mapping to different emphasis within the 15 grammar operations. All four share the same underlying data model (nodes, ops, event graph) but present radically different interfaces.

| Mode | Replaces | Primary Grammar Ops | Core Pattern |
|------|----------|-------------------|--------------|
| **Chat** | Slack, Messenger, iMessage | Emit, Respond, Channel, Subscribe | Real-time 1:1 and group messaging |
| **Rooms** | Discord, IRC | Channel, Subscribe, Delegate, Consent | Persistent community spaces with channels |
| **Square** | Twitter/X, Bluesky | Emit, Propagate, Endorse, Annotate | Public broadcast with engagement mechanics |
| **Forum** | Reddit, Discourse | Emit, Respond, Endorse, Channel, Merge | Structured threaded discussion with quality signals |

The sidebar restructures Social as:

```
Social
├── Chat       (DMs + group conversations)
├── Rooms      (community channels)
├── Square     (public feed)
└── Forum      (discussions)
```

### Shared Infrastructure

All four modes use these Code Graph primitives as shared infrastructure:

```
Entity(Message, properties: [body, author_id, reply_to, reactions, edited_at])
Entity(Conversation, properties: [title, kind, participants, last_message_at])
Entity(Channel, properties: [name, topic, space_id, kind, position])
Entity(Post, properties: [title, body, author_id, engagement])

State(MessageState, values: [active, edited, deleted])
State(ThreadState, values: [open, locked, archived])

Relation(Message -> Conversation, type: belongs_to)
Relation(Message -> Message, type: reply_to)
Relation(User -> Conversation, type: member_of)

Event(message.created, message.edited, message.deleted, message.reacted)
Event(conversation.created, member.joined, member.left)

Subscribe(Query(Message, filter: conversation.current), on_change: View.append)
```

---

## Mode 1: Chat

**Goal:** Beat Slack and iMessage in speed, clarity, and agent integration.

### What makes Slack/iMessage great (and what we must match)

1. **Message grouping** — consecutive messages from same author within 5min collapse into a single visual block (reduces clutter ~40%)
2. **Threading** — reply to any message to create a sub-conversation without polluting the main flow
3. **Reactions** — emoji on any message as lightweight acknowledgment (avoids "thanks!" noise)
4. **Typing indicators** — see when someone is composing
5. **Read receipts** — know your message was seen
6. **Mentions** — @person for directed attention, @here for present members
7. **Rich content** — markdown, code blocks, embeds, images, files
8. **Search** — find any message by content, sender, date, channel
9. **Keyboard-first** — up-arrow to edit last, /commands, Cmd+K navigation
10. **Huddles** — spontaneous voice/video (Slack)
11. **Slash commands** — /remind, /poll, extensible

### What we add that they can't

- **Endorse** on messages — not just react with emoji, but formally endorse a statement (recorded on event graph, builds author reputation)
- **Delegate** — transfer conversation ownership/responsibility with audit trail
- **Consent** — structured agreement on decisions within chat (not just informal "sounds good")
- **Annotate** — structured commentary that adds metadata, not just inline text
- **Agent as peer** — AI participates naturally with identity, not as a bot add-on

### Code Graph Composition

```
Chat = View(name: ChatApp,
  layout: Layout(split, ratio: [1, 3]),

  // Left panel: conversation list
  sidebar: Layout(stack, [
    Input(search, type: text, placeholder: "Search conversations..."),
    Action(label: "+", command: Navigation(modal: Form(create_conversation))),
    List(Query(Conversation,
        filter: { participants: contains(current_user) },
        sort: last_message_at.desc),
      template: ConversationCard(
        avatar: Avatar(conversation.other_participant),
        title: Display(conversation.title),
        preview: Display(conversation.last_message.body, truncate: 60),
        time: Recency(conversation.last_message_at),
        unread: Display(conversation.unread_count, style: badge),
        presence: Presence(conversation.other_participant)
      ))
  ]),

  // Right panel: active conversation
  main: Layout(stack, [
    // Header
    Layout(row, justify: space-between, [
      Layout(row, [
        Avatar(conversation.participants, size: sm),
        Display(conversation.title, style: heading),
        Presence(conversation.participants)
      ]),
      Layout(row, [
        Action(label: "search", command: Toggle(message_search)),
        Action(label: "info", command: Toggle(conversation_info))
      ])
    ]),

    // Messages
    List(Query(Message,
        filter: { conversation_id: current },
        sort: created_at.asc),
      template: MessageBubble,
      subscribe: Subscribe(on_change: append),
      pagination: Pagination(type: load_more, direction: up, page_size: 50),
      grouping: GroupBy(author_id, window: 5m)),

    // Typing indicator
    Liveness(collaborators: Query(Session, filter: { typing_in: conversation }),
      display: Display("{names} typing...", style: caption)),

    // Compose
    Layout(row, [
      if reply_to {
        Layout(row, [
          Display(reply_to.preview, style: caption),
          Action(label: "×", command: Command(clear_reply))
        ])
      },
      Input(body, type: rich_text,
        placeholder: "Message...",
        shortcuts: {
          enter: submit,
          shift_enter: newline,
          up_arrow: edit_last,
          escape: clear
        }),
      Action(label: "Send", command: Sequence([
        Command(create, Entity(Message, {
          body: input.body,
          conversation_id: current,
          reply_to: reply_to.id
        })),
        Feedback(optimistic, display: message_preview)
      ]))
    ])
  ])
)
```

### Message Component

```
MessageBubble = Condition(
  // Grouped: compact render
  msg.grouped, then:
    Layout(row, class: "ml-10", [
      Display(msg.body, format: markdown),
      Selection(mode: hover, actions: [
        Action(label: "react", command: Navigation(popover: EmojiPicker)),
        Action(label: "reply", command: Command(set_reply_to, msg)),
        Action(label: "endorse", command: Command(endorse, msg)),
        Action(label: "...", command: Navigation(popover: MessageMenu))
      ])
    ]),

  // Full: avatar + name + time
  else:
    Layout(row, align: top, [
      Avatar(msg.author, size: md),
      Layout(stack, [
        Layout(row, gap: sm, [
          Display(msg.author.name, style: bold),
          Display(msg.created_at, style: relative_time)
        ]),
        Display(msg.body, format: markdown),
        if msg.reactions.length > 0 {
          Layout(row, wrap: true, gap: xs,
            Loop(msg.reactions, each: r ->
              Action(label: Display(r.emoji + " " + r.count),
                style: if r.includes(current_user) then "active",
                command: Command(toggle_reaction, msg, r.emoji))))
        },
        if msg.reply_to {
          Display(msg.reply_to.preview, style: quote)
        }
      ]),
      Selection(mode: hover, actions: [...])
    ])
)
```

### Differentiator: Endorse in Chat

```
// When user endorses a message, it's recorded as a grammar op on the event graph
Trigger(on: Message.endorsed,
  do: Sequence([
    Command(Event.emit(endorse, {
      target: message,
      actor: current_user,
      evidence: message.body
    })),
    // Builds author's reputation score
    Command(update, Author.reputation, increment: endorsement_weight),
    // Visible as a special reaction
    Display(endorsement_count, style: badge(brand))
  ]))
```

### Differentiator: Consent in Chat

```
// Any message can be promoted to a decision point
ConsentRequest = Form(
  trigger: Action(label: "Request consent", on: Message),
  fields: [
    Display(message.body, style: quote),
    Input(question, type: text, default: "Do you agree?"),
    Input(participants, type: entity_picker, default: conversation.participants)
  ],
  submit: Command(create, Entity(ConsentRequest, {
    message_id: message.id,
    question: input.question,
    required: input.participants
  }))
)

// Renders inline as a structured decision
ConsentCard = Layout(stack, [
  Display(consent.question, style: heading),
  Display(consent.message.body, style: quote),
  Layout(row, Loop(consent.participants, each: p ->
    Layout(row, [
      Avatar(p, size: sm),
      Condition(
        consent.voted(p), then: Display(consent.vote(p), style: badge),
        else: Display("pending", style: muted)
      )
    ])
  )),
  if !consent.voted(current_user) {
    Layout(row, [
      Action(label: "Agree", command: Command(consent, { vote: yes })),
      Action(label: "Disagree", command: Command(dissent, { vote: no, reason: Input(text) }))
    ])
  },
  Display(consent.status, style: status_badge) // "2/3 agreed", "Consensus reached"
])
```

### Build Priority (Chat)

| # | Feature | Effort | Impact | Grammar Op |
|---|---------|--------|--------|-----------|
| 1 | Reactions (emoji on messages) | S | High | Acknowledge |
| 2 | Reply-to linkage (fix existing) | S | High | Respond |
| 3 | Message editing | S | Medium | Extend |
| 4 | Message deletion (soft) | S | Medium | Retract |
| 5 | Unread count per conversation | S | High | Subscribe |
| 6 | Typing indicator | M | Medium | Liveness |
| 7 | Message search | M | High | Search |
| 8 | DM vs group distinction | S | Medium | Channel |
| 9 | Endorse on messages | M | High | Endorse |
| 10 | Consent requests inline | L | Differentiator | Consent |
| 11 | Conversation archive/mute | S | Medium | Subscribe |
| 12 | File/image attachments | L | Medium | — |

---

## Mode 2: Rooms

**Goal:** Beat Discord in community organization. Match the "always-on space" feel.

### What makes Discord great (and what we must match)

1. **Server → Category → Channel hierarchy** — organize a community's communication into logical groups
2. **Channel types** — text, voice, forum, stage, announcement
3. **Persistent channels** — you join a room, it's always there, history preserved
4. **Roles + permissions** — fine-grained per-channel, per-role access control
5. **Rich presence** — see who's online, what they're doing
6. **Message pinning** — pin important messages to channel
7. **Slowmode** — rate limit messages (anti-spam, focus discussion)
8. **Forum channels** — structured posts with tags, solved markers
9. **Thread auto-archive** — threads go quiet after inactivity, declutter
10. **Server discovery** — find communities by category

### What we add that Discord can't

- **Endorse** — formally vouch for members (builds verifiable reputation, not just role assignment)
- **Delegate** — transfer channel/category ownership with full audit trail
- **Consent** — community decisions require explicit consent, not just moderator fiat
- **Merge** — combine duplicate channels/threads with provenance preserved
- **Annotate** — add structured metadata to any message (not just pin/react)
- **Agent members** — AI agents as first-class community participants with visible identity

### Code Graph Composition

```
Rooms = View(name: RoomsApp,
  layout: Layout(split, ratio: [1, 1, 4]),

  // Left: Space list
  spaces: Layout(stack, [
    Loop(Query(Space, filter: { member: current_user }),
      template: SpaceIcon(
        avatar: Avatar(space, size: lg),
        unread: Display(space.unread_count, style: dot),
        name: Display(space.name, style: tooltip)
      ))
  ]),

  // Middle: Channel list for active space
  channels: Layout(stack, [
    Display(space.name, style: heading),
    Loop(Query(Category, filter: { space_id: current_space }, sort: position),
      template: Layout(stack, [
        Display(category.name, style: subheading),
        Loop(Query(Channel, filter: { category_id: category.id }, sort: position),
          template: ChannelItem(
            icon: Display(channel.kind, style: icon), // #, voice, forum
            name: Display(channel.name),
            unread: Display(channel.unread_count, style: badge),
            active: Presence(scope: channel)
          ))
      ])),
    // Members panel (collapsible)
    Layout(stack, [
      Display("Members", style: subheading),
      Loop(Query(Member, filter: { space_id: current_space, online: true }),
        template: Layout(row, [
          Avatar(member, size: sm),
          Display(member.name),
          Presence(member, display: dot)
        ]))
    ])
  ]),

  // Right: Active channel content
  main: Condition(
    channel.kind == "text", then: TextChannel(channel),
    channel.kind == "forum", then: ForumChannel(channel),
    channel.kind == "voice", then: VoiceChannel(channel)
  )
)

TextChannel = Layout(stack, [
  // Channel header
  Layout(row, justify: space-between, [
    Layout(row, [
      Display("#", style: icon),
      Display(channel.name, style: heading),
      Display(channel.topic, style: muted)
    ]),
    Layout(row, [
      Action(label: "threads", command: Toggle(thread_panel)),
      Action(label: "pins", command: Toggle(pins_panel)),
      Action(label: "members", command: Toggle(members_panel)),
      Input(search, type: text, placeholder: "Search")
    ])
  ]),

  // Messages (same MessageBubble as Chat, but with channel context)
  List(Query(Message,
      filter: { channel_id: current },
      sort: created_at.asc),
    template: MessageBubble,
    subscribe: Subscribe(on_change: append),
    pagination: Pagination(type: load_more, direction: up, page_size: 50),
    grouping: GroupBy(author_id, window: 5m)),

  // Compose (same as Chat)
  ComposeBar(channel)
])

ForumChannel = Layout(stack, [
  // Header with create button
  Layout(row, justify: space-between, [
    Display(channel.name, style: heading),
    Action(label: "New Post", command: Navigation(modal: Form(create_forum_post)))
  ]),

  // Post list
  List(Query(ForumPost,
      filter: { channel_id: current },
      sort: Condition(
        sort_mode == "hot", then: hotness_score.desc,
        sort_mode == "new", then: created_at.desc,
        sort_mode == "top", then: endorsement_count.desc
      )),
    template: ForumPostCard(
      title: Display(post.title),
      author: Avatar(post.author),
      preview: Display(post.body, truncate: 200),
      tags: Loop(post.tags, each: Display(tag, style: badge)),
      replies: Display(post.reply_count),
      endorsements: Display(post.endorsement_count),
      status: Display(post.state, style: badge) // open, solved, locked
    ))
])
```

### Differentiator: Community Governance via Consent

```
// Channel creation requires community consent (configurable)
Trigger(on: Channel.create_request,
  do: Condition(
    space.governance == "consent",
    then: Sequence([
      Command(create, Entity(Proposal, {
        type: "create_channel",
        title: "Create #" + request.name,
        description: request.reason,
        required_consent: space.quorum
      })),
      // Channel only created when consent threshold met
      Trigger(on: Proposal.passed,
        do: Command(create, Entity(Channel, request.config)))
    ]),
    else: // Owner creates directly
      Command(create, Entity(Channel, request.config))
  ))
```

### Build Priority (Rooms)

| # | Feature | Effort | Impact | Grammar Op |
|---|---------|--------|--------|-----------|
| 1 | Channel model (text channels within spaces) | M | Critical | Channel |
| 2 | Channel list sidebar | M | Critical | Subscribe |
| 3 | Channel-level messages | S | Critical | Emit, Respond |
| 4 | Categories (channel groups) | S | High | Channel |
| 5 | Channel topics | S | Medium | Annotate |
| 6 | Roles + permissions | L | High | Delegate |
| 7 | Member online status | M | Medium | Presence |
| 8 | Channel pinned messages | S | Medium | — |
| 9 | Forum-type channels | M | High | Emit, Endorse |
| 10 | Thread sidebar | M | Medium | Respond |

---

## Mode 3: Square

**Goal:** Beat Twitter in public discourse quality. Same reach, better signal-to-noise.

### What makes Twitter great (and what we must match)

1. **Compose + publish in seconds** — minimal friction to emit thought
2. **The feed** — algorithmic or chronological, shows content from people you follow
3. **Engagement chain** — like, repost, quote, reply form a rich interaction graph
4. **Quote posts** — respond to content while adding your own context (Derive)
5. **Threading** — multi-post threads for long-form (self-Extend)
6. **Bookmarks** — save for later without signaling publicly
7. **Lists** — curated feeds by topic or people
8. **Trending** — what the community is talking about right now
9. **Profile** — identity, bio, pinned post, follower/following counts
10. **Notifications** — mentions, likes, reposts, follows, quote posts

### What we add that Twitter can't

- **Endorse** vs Like — endorsement is a verifiable reputation signal recorded on the event graph, not a throwaway metric. Endorsing means "I stand behind this statement."
- **Consent** — polls with formal consent semantics, not just engagement bait
- **Annotate** — structured fact-checks, context, corrections attached to posts (community notes but as a primitive)
- **Merge** — combine duplicate conversations/threads preserving provenance
- **Sever** — formal relationship severance with recorded reason (not just block/mute)
- **Audit trail** — every interaction signed, hash-chained, causally linked. No shadow bans, no invisible algorithmic suppression.

### Code Graph Composition

```
Square = View(name: SquareApp,
  layout: Layout(stack),

  // Compose
  Form(command: Command(emit, Entity(Post)), fields: [
    Avatar(current_user, size: md),
    Input(body, type: rich_text, max_length: 1000,
      placeholder: "What's on your mind?"),
    Layout(row, justify: space-between, [
      Layout(row, [
        Action(label: "image", command: Input(file, type: image)),
        Action(label: "poll", command: Toggle(poll_form)),
        Action(label: "consent", command: Toggle(consent_form))
      ]),
      Action(label: "Post", command: submit,
        condition: Constraint(body.length > 0))
    ])
  ]),

  // Feed
  Layout(row, [
    Action(label: "Following", style: if mode == "following" then "active"),
    Action(label: "For You", style: if mode == "foryou" then "active"),
    Action(label: "Trending", style: if mode == "trending" then "active")
  ]),

  List(Query(Post,
      filter: Condition(
        mode == "following", then: { author: current_user.following },
        mode == "foryou", then: { score: algorithmic },
        mode == "trending", then: { trending: true }
      ),
      sort: Condition(
        mode == "trending", then: hotness.desc,
        else: created_at.desc
      )),
    template: PostCard,
    subscribe: Subscribe(on_change: prepend),
    pagination: Pagination(type: infinite_scroll, page_size: 20))
)

PostCard = Layout(stack, [
  // Repost header (if reposted)
  if post.reposted_by {
    Layout(row, [
      Display("reposted", style: icon),
      Display(post.reposted_by.name, style: caption)
    ])
  },

  Layout(row, align: top, gap: md, [
    Avatar(post.author, size: md),
    Layout(stack, flex: 1, [
      Layout(row, [
        Display(post.author.name, style: bold),
        Display(post.author.handle, style: muted),
        Display("·"),
        Recency(post.created_at)
      ]),

      // Quote context (if quoting another post)
      if post.quote_of {
        Layout(stack, class: "border rounded p-3 mt-1", [
          Layout(row, [
            Avatar(post.quote_of.author, size: xs),
            Display(post.quote_of.author.name, style: bold_sm),
            Recency(post.quote_of.created_at, style: caption)
          ]),
          Display(post.quote_of.body, truncate: 100)
        ])
      },

      Display(post.body, format: markdown),

      // Engagement bar
      Layout(row, justify: space-between, class: "mt-2", [
        Action(label: Display(post.reply_count),
          icon: "reply",
          command: Navigation(route: post_detail)),

        Action(label: Display(post.repost_count),
          icon: "repost",
          command: Command(propagate, post)),

        Action(label: Display(post.endorsement_count),
          icon: "endorse",
          style: if post.endorsed_by(current_user) then "active",
          command: Command(endorse, post)),

        Action(label: "bookmark",
          icon: "bookmark",
          style: if post.bookmarked then "active",
          command: Command(bookmark, post)),

        Action(label: "...",
          command: Navigation(popover: PostMenu(post)))
      ]),

      // Annotations (community notes / fact checks)
      if post.annotations.length > 0 {
        Layout(stack, class: "bg-surface rounded p-3 mt-2", [
          Display("Community context", style: subheading_sm),
          Loop(post.annotations, each: a ->
            Layout(stack, [
              Display(a.body, format: markdown),
              Layout(row, [
                Display("Helpful?"),
                Action(label: "Yes (" + a.helpful_count + ")",
                  command: Command(endorse, a)),
                Action(label: "No",
                  command: Command(dissent, a))
              ])
            ]))
        ])
      }
    ])
  ])
])
```

### Differentiator: Annotate as Primitive

```
// Any post can be annotated with structured context
// Annotations are community-contributed, quality-ranked by endorsements
Annotation = Sequence([
  Form(fields: [
    Display(post.body, style: quote),
    Input(annotation_type, type: select, options: [
      "context",      // adds missing context
      "correction",   // factual correction with evidence
      "source",       // links to primary source
      "perspective",  // alternative viewpoint
    ]),
    Input(body, type: rich_text),
    Input(evidence, type: url, optional: true)
  ]),
  Command(annotate, {
    target: post,
    type: annotation_type,
    body: input.body,
    evidence: input.evidence
  })
])

// Annotations surface when they cross an endorsement threshold
Trigger(on: Annotation.endorsement_count > threshold,
  do: Command(update, Post.annotations, append: annotation))
```

### Build Priority (Square)

| # | Feature | Effort | Impact | Grammar Op |
|---|---------|--------|--------|-----------|
| 1 | Public post feed (existing Feed, enhanced) | M | Critical | Emit |
| 2 | Endorse (replace like) | S | Critical | Endorse |
| 3 | Repost / Propagate | S | High | Propagate |
| 4 | Quote post (post-with-reference) | M | High | Derive |
| 5 | Reply threads on posts | S | High | Respond |
| 6 | Follow users | M | High | Subscribe |
| 7 | Engagement counts on cards | S | High | — |
| 8 | Trending / hot algorithm | M | Medium | — |
| 9 | Bookmarks | S | Medium | — |
| 10 | Annotations (community context) | L | Differentiator | Annotate |
| 11 | User lists / curated feeds | M | Medium | Channel |

---

## Mode 4: Forum

**Goal:** Beat Reddit in discussion quality. Same depth, better governance.

### What makes Reddit great (and what we must match)

1. **Subreddits** — topic-specific communities with their own rules, moderators, culture
2. **Upvote/downvote** — community quality signal drives content ranking
3. **Comment threading** — infinite nesting for deep discussions
4. **Sorting** — hot, best, new, top, controversial (each optimizes for different signal)
5. **Post types** — text, link, image, video, poll
6. **Flair** — categorize posts and users
7. **Moderator tools** — automod, rules, post removal, user bans
8. **Cross-posting** — share between communities
9. **Wikis** — collaborative knowledge bases per community
10. **Awards** — community signaling of exceptional content

### What we add that Reddit can't

- **Endorse** vs Upvote — endorsement is identity-linked, reputation-building, auditable. Not anonymous, not gameable.
- **Consent** — community decisions use formal consent, not mod fiat
- **Delegate** — moderators have explicit delegated authority from community, revocable
- **Merge** — duplicate threads merged preserving both conversations
- **Annotate** — structured responses (correction, source, context) distinct from discussion
- **Fork** — split a discussion thread when it's covering multiple topics

### Code Graph Composition

```
Forum = View(name: ForumApp,
  layout: Layout(stack),

  // Sort controls
  Layout(row, [
    Action(label: "Hot", command: sort("hot")),
    Action(label: "New", command: sort("new")),
    Action(label: "Top", command: sort("top")),
    Action(label: "Discussed", command: sort("discussed"))
  ]),

  // Create
  Action(label: "New Discussion",
    command: Navigation(modal: Form(create_discussion, fields: [
      Input(title, type: text, required: true),
      Input(body, type: rich_text),
      Input(tags, type: multi_select, options: Query(Tag, filter: space.current)),
      Input(type, type: select, options: ["discussion", "question", "proposal"])
    ]))),

  // Thread list
  List(Query(Thread,
      filter: { space_id: current, kind: "discussion" },
      sort: current_sort),
    template: ThreadCard(
      // Vote column
      Layout(stack, align: center, [
        Action(label: "▲", command: Command(endorse, thread),
          style: if thread.endorsed_by(current_user) then "active"),
        Display(thread.score, style: bold),
        Action(label: "▼", command: Command(dissent, thread),
          style: if thread.dissented_by(current_user) then "active")
      ]),

      // Content
      Layout(stack, [
        // Tags
        Layout(row, gap: xs,
          Loop(thread.tags, each: Display(tag, style: badge))),

        // Title + meta
        Display(thread.title, style: heading),
        Layout(row, gap: sm, [
          Avatar(thread.author, size: xs),
          Display(thread.author.name),
          Display("·"),
          Recency(thread.created_at),
          Display("·"),
          Display(thread.reply_count + " replies"),
          if thread.state == "solved" {
            Display("Solved", style: badge(success))
          }
        ]),

        // Preview
        Display(thread.body, truncate: 200)
      ])
    ),
    pagination: Pagination(type: infinite_scroll, page_size: 25))
)

// Thread detail with nested comments
ThreadDetail = Layout(stack, [
  // Original post
  Layout(stack, [
    Display(thread.title, style: display),
    Layout(row, [
      Avatar(thread.author),
      Display(thread.author.name, style: bold),
      Recency(thread.created_at)
    ]),
    Display(thread.body, format: markdown),
    // Engagement bar (endorse, reply, share, annotate)
    EngagementBar(thread)
  ]),

  // Sort comments
  Layout(row, [
    Display("Comments (" + thread.reply_count + ")"),
    Action(label: "Best", command: sort("best")),
    Action(label: "New", command: sort("new")),
    Action(label: "Top", command: sort("top"))
  ]),

  // Nested comments (recursive)
  CommentTree(thread)
])

CommentTree = Loop(Query(Comment,
    filter: { parent_id: thread.id, depth: 0 },
    sort: current_sort),
  template: CommentNode)

CommentNode = Layout(stack, class: "ml-" + depth * 4, [
  Layout(row, align: top, [
    // Vote
    Layout(stack, align: center, [
      Action(label: "▲", command: Command(endorse, comment)),
      Display(comment.score),
      Action(label: "▼", command: Command(dissent, comment))
    ]),

    Layout(stack, [
      Layout(row, [
        Avatar(comment.author, size: xs),
        Display(comment.author.name, style: bold),
        Recency(comment.created_at),
        if comment.edited_at {
          Display("(edited)", style: muted)
        }
      ]),
      Display(comment.body, format: markdown),

      // Actions
      Layout(row, [
        Action(label: "Reply", command: Toggle(reply_form)),
        Action(label: "Endorse", command: Command(endorse, comment)),
        Action(label: "Annotate", command: Toggle(annotate_form)),
        if comment.children_count > collapse_threshold {
          Action(label: "Collapse (" + comment.children_count + ")",
            command: Toggle(collapsed))
        }
      ]),

      // Nested children (recursive)
      if !collapsed && comment.children_count > 0 {
        Loop(Query(Comment, filter: { parent_id: comment.id }, sort: best),
          template: CommentNode(depth: depth + 1))
      }
    ])
  ])
])
```

### Differentiator: Merge and Fork

```
// Merge: combine duplicate threads
Merge = Sequence([
  // Show both threads side by side
  Layout(split, [
    Display(thread_a, style: preview),
    Display(thread_b, style: preview)
  ]),
  Confirmation(
    message: "Merge these threads? Comments from both will be combined.",
    consequence: Consequence Preview(
      impact: [
        { label: "comments combined", count: thread_a.count + thread_b.count },
        { label: "participants merged", count: unique(participants) }
      ]
    )
  ),
  Command(merge, { source: thread_b, target: thread_a }),
  // Both threads' histories preserved on event graph
  // Redirect from thread_b to thread_a
  Event(merge, { source: thread_b, target: thread_a, reason: "duplicate" })
])

// Fork: split a thread when discussion has diverged
Fork = Sequence([
  Selection(scope: CommentTree, mode: branch,
    prompt: "Select the comment branch to split off"),
  Form(fields: [
    Input(title, type: text, placeholder: "New thread title"),
    Input(reason, type: text, placeholder: "Why split?")
  ]),
  Command(fork, {
    source: thread,
    branch: selected_comments,
    title: input.title
  }),
  Event(fork, { source: thread, target: new_thread, reason: input.reason })
])
```

### Build Priority (Forum)

| # | Feature | Effort | Impact | Grammar Op |
|---|---------|--------|--------|-----------|
| 1 | Discussion threads (enhance existing) | M | Critical | Emit |
| 2 | Nested comment trees | M | Critical | Respond |
| 3 | Endorse/Dissent voting | S | Critical | Endorse |
| 4 | Score-based sorting (hot, top, new) | M | High | — |
| 5 | Post tags / flair | S | Medium | Annotate |
| 6 | Solved/resolved state | S | Medium | — |
| 7 | Comment collapse | S | Medium | — |
| 8 | Merge duplicate threads | L | Differentiator | Merge |
| 9 | Fork diverged discussions | L | Differentiator | Fork |
| 10 | Cross-post between spaces | M | Medium | Propagate |

---

## Grammar Operation Coverage Matrix

Every grammar operation should be exercised by at least one mode:

| Operation | Chat | Rooms | Square | Forum | Description |
|-----------|------|-------|--------|-------|-------------|
| **Emit** | Send message | Post in channel | Create post | Create thread | Produce content |
| **Respond** | Reply to message | Reply in channel | Reply to post | Nested comment | Direct response |
| **Derive** | — | — | Quote post | — | New content from existing |
| **Extend** | Edit message | Edit message | Edit post | Edit comment | Modify own content |
| **Retract** | Delete message | Delete message | Delete post | Delete comment | Remove content |
| **Annotate** | — | Channel topic | Community context | Fact-check | Add structured metadata |
| **Acknowledge** | Reaction (emoji) | Reaction | — | — | Lightweight signal |
| **Propagate** | Forward | Cross-post | Repost | Cross-post | Share to new audience |
| **Endorse** | Endorse message | Endorse post | Endorse post | Upvote | Reputation signal |
| **Subscribe** | Join conversation | Join channel | Follow user | Join space | Opt into updates |
| **Channel** | Conversation | Channel | List | Space/tag | Content routing |
| **Delegate** | Transfer ownership | Assign moderator | — | Assign moderator | Authority transfer |
| **Consent** | Decision request | Community vote | Poll consent | Governance proposal | Structured agreement |
| **Sever** | Leave/block | Leave/kick | Unfollow/block | Leave/ban | End relationship |
| **Merge** | — | — | — | Merge threads | Combine content |

**Coverage: 15/15 operations exercised across the four modes.** No operation is left unused.

---

## Database Changes Required

### New columns on `nodes`

```sql
ALTER TABLE nodes ADD COLUMN reply_to_id TEXT;          -- for message reply threading
ALTER TABLE nodes ADD COLUMN edited_at TIMESTAMPTZ;     -- track edits
ALTER TABLE nodes ADD COLUMN score INTEGER DEFAULT 0;   -- endorsement - dissent
ALTER TABLE nodes ADD COLUMN channel_id TEXT;            -- for Rooms mode
```

### New table: `reactions`

```sql
CREATE TABLE reactions (
  node_id TEXT NOT NULL,
  user_id TEXT NOT NULL,
  emoji TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (node_id, user_id, emoji)
);
```

### New table: `follows`

```sql
CREATE TABLE follows (
  follower_id TEXT NOT NULL,
  following_id TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (follower_id, following_id)
);
```

### New table: `bookmarks`

```sql
CREATE TABLE bookmarks (
  user_id TEXT NOT NULL,
  node_id TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (user_id, node_id)
);
```

### New table: `channels`

```sql
CREATE TABLE channels (
  id TEXT PRIMARY KEY,
  space_id TEXT NOT NULL,
  name TEXT NOT NULL,
  topic TEXT DEFAULT '',
  kind TEXT NOT NULL DEFAULT 'text',  -- text, forum, announcement
  category TEXT DEFAULT '',
  position INTEGER DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### New table: `read_state`

```sql
CREATE TABLE read_state (
  user_id TEXT NOT NULL,
  conversation_id TEXT NOT NULL,
  last_read_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (user_id, conversation_id)
);
```

---

## Implementation Order

**Phase 1: Chat Foundation** (iters 183-188)
Core messaging improvements that benefit ALL modes.
1. Reactions (emoji)
2. Reply-to linkage (fix)
3. Message edit/delete
4. Unread counts
5. DM vs group distinction
6. Message search

**Phase 2: Square** (iters 189-194)
Public discourse builds on the existing Feed lens.
1. Endorse on posts (replace like)
2. Repost/propagate
3. Quote post
4. Follow users
5. Engagement counts
6. Trending algorithm

**Phase 3: Forum** (iters 195-200)
Structured discussion builds on existing Threads lens.
1. Endorse/dissent voting on threads
2. Nested comment trees
3. Score-based sorting
4. Post tags
5. Solved state
6. Comment collapse

**Phase 4: Rooms** (iters 201-208)
Community channels require the most new infrastructure.
1. Channel model + table
2. Channel list UI
3. Channel messages
4. Categories
5. Channel-level permissions
6. Forum-type channels
7. Roles
8. Member presence

**Phase 5: Differentiators** (iters 209-215)
Our unique grammar ops as features.
1. Consent requests in chat
2. Annotations on posts
3. Merge threads
4. Fork discussions
5. Delegate moderation
6. Cross-space posting
7. Audit trail UI

---

## Success Criteria

For each mode, "competitive" means a user who currently uses the incumbent would find our version **at least as good** for their primary workflow, plus **meaningfully better** in at least one dimension they care about.

| Mode | Minimum Bar | Differentiator |
|------|------------|----------------|
| Chat | Fast, grouped, reactions, search, threads | Consent, Endorse, Agent peers |
| Rooms | Persistent channels, categories, permissions | Governance, Delegate, Merge |
| Square | Fast compose, rich feed, engagement | Annotations, Endorse reputation, Audit |
| Forum | Nested threads, voting, sorting | Merge, Fork, formal Consent |

---

*This spec is the build blueprint. Every iteration references it. Every feature traces back to a Code Graph primitive and a grammar operation. Ship from spec, not from intuition.*
